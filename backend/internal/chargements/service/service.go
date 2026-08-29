package chargementservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/KochKevin/effective-spoon-v2/internal/chargements"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v86"
)

type Repo interface {
	SetStripeLastEventId(ctx context.Context, tx *sql.Tx, lastEventId string) error
	GetStripeLastEventId(ctx context.Context) (string, error)

	CreateBalanceChargementIntent(ctx context.Context, tx *sql.Tx, chargementIntent chargements.ChargementIntent) error
	GetBalanceChargementIntent(ctx context.Context, tx *sql.Tx, id uuid.UUID) (chargements.ChargementIntent, error)
	SaveBalanceChargementIntent(ctx context.Context, tx *sql.Tx, chargementIntent chargements.ChargementIntent) error
}

type UserRepo interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction users.Transaction) (users.Transaction, error)
}

type Service struct {
	Repo                         Repo
	BalanceChargementIntentCache BalanceChargementIntentCache
	UserRepo                     UserRepo
	Txm                          infrastructure.TxManager
	StripeClient                 *stripe.Client
	PushService                  PushService
}

func (s *Service) CreateChargementIntent(ctx context.Context, userId uuid.UUID, amount float32) (chargements.ChargementIntent, error) {

	chargementIntent, err := chargements.NewChargementIntent(money.MoneyFrom(int(amount*100)), userId)
	if err != nil {
		return chargements.ChargementIntent{}, fmt.Errorf("error creating chargement intent: %w", err)
	}

	chargementIntent, err = s.createStripePaymentLink(ctx, chargementIntent)
	if err != nil {
		return chargements.ChargementIntent{}, fmt.Errorf("error creating stripe payment link: %w", err)
	}

	err = s.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		err := s.Repo.CreateBalanceChargementIntent(context.Background(), tx, chargementIntent)

		if err != nil {
			return fmt.Errorf("error creating chargement intent on persitent volume: %w", err)
		}

		//Set Current
		go s.BalanceChargementIntentCache.SetBalanceChargementIntentId(chargementIntent.Id)

		return nil
	})
	if err != nil {
		return chargements.ChargementIntent{}, fmt.Errorf("error in transaction: %w", err)
	}

	return chargementIntent, nil

}

const MetadataBalanceChargementIntentId = "balance_chargement_intent_id"
const PaymentLinkCustomMessage = "🎉 Danke für das Aufladen deines Kontos. Du kannst den Browser nun schließen. Die Verarbeitung deiner Aufladung kann bis zu 10 Sekunden dauern 💸"
const PaymentLinkChargementIntentProductName = "%.2f€ Getränkekasse Guthabenaufladung"

func (s *Service) createStripePaymentLink(ctx context.Context, chargementIntent chargements.ChargementIntent) (chargements.ChargementIntent, error) {

	metadata := map[string]string{
		MetadataBalanceChargementIntentId: chargementIntent.Id.String(),
	}

	params := &stripe.PaymentLinkCreateParams{
		Metadata: metadata,
		AfterCompletion: &stripe.PaymentLinkCreateAfterCompletionParams{
			Type: stripe.String("hosted_confirmation"),
			HostedConfirmation: &stripe.PaymentLinkCreateAfterCompletionHostedConfirmationParams{
				CustomMessage: stripe.String(PaymentLinkCustomMessage),
			},
		},
		LineItems: []*stripe.PaymentLinkCreateLineItemParams{
			&stripe.PaymentLinkCreateLineItemParams{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.PaymentLinkCreateLineItemPriceDataParams{
					Currency: stripe.String("EUR"),
					ProductData: &stripe.PaymentLinkCreateLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf(PaymentLinkChargementIntentProductName, chargementIntent.Amount.GetAsEuro())),
					},
					UnitAmount: stripe.Int64(int64(chargementIntent.Amount.GetAsCents())),
				},
			},
		},
	}

	result, err := s.StripeClient.V1PaymentLinks.Create(ctx, params)
	if err != nil {
		return chargements.ChargementIntent{}, fmt.Errorf("error calling create stripe payment link api: %v", err)
	}

	chargementIntent.AddStripePaymentLinkAndId(result.URL, result.ID)

	return chargementIntent, nil

}

func (s *Service) deactivateStripePaymentLink(ctx context.Context, paymentLinkId string) error {

	params := &stripe.PaymentLinkUpdateParams{}
	params.Active = stripe.Bool(false)
	_, err := s.StripeClient.V1PaymentLinks.Update(ctx, paymentLinkId, params)
	if err != nil {
		return fmt.Errorf("error while deactivating stripe payment link: %w", err)
	}

	return nil

}

func (s *Service) CancelCurrentChargementIntent(ctx context.Context) error {

	err := s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		chargementIntentId := s.BalanceChargementIntentCache.GetBalanceChargementIntentId()

		if chargementIntentId == uuid.Nil {
			return errors.New("no current chargement active")
		}

		chargementIntent, err := s.Repo.GetBalanceChargementIntent(ctx, tx, chargementIntentId)

		if err != nil {
			return fmt.Errorf("error getting chargement intent from persitent volume: %w", err)
		}

		err = s.deactivateStripePaymentLink(ctx, chargementIntent.PaymentLinkId)
		if err != nil {
			return fmt.Errorf("an error occured while deactivting an stripe payment link: %w", err)
		}

		chargementIntent.Status = chargements.StatusCanceld

		err = s.Repo.SaveBalanceChargementIntent(ctx, tx, chargementIntent)
		if err != nil {
			return fmt.Errorf("error saving chargement intent to persitent volume: %w", err)
		}

		s.BalanceChargementIntentCache.ClearBalanceChargementIntentId()

		return nil
	})

	if err != nil {
		return fmt.Errorf("error in cancel current chargement intent transaction: %w", err)
	}

	return nil
}
