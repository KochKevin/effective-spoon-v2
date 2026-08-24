package chargementservice

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/chargements"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v86"
)

const MetadataBalanceChargementIntentId = "balance_chargemen_intent_id"

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
}

func (s *Service) CreatePaymentLink() (string, error) {

	var chargemnet chargements.ChargementIntent

	err := s.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		var err error
		chargemnet, err = chargements.NewChargementIntent(money.MoneyFrom(500), uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		if err != nil {
			return err
		}

		s.Repo.CreateBalanceChargementIntent(context.Background(), tx, chargemnet)

		return nil
	})
	if err != nil {
		slog.Error("error in transaction", err)
		return "", err
	}

	metadata := map[string]string{
		MetadataBalanceChargementIntentId: chargemnet.Id.String(),
	}

	params := &stripe.PaymentLinkCreateParams{
		Metadata: metadata,
		AfterCompletion: &stripe.PaymentLinkCreateAfterCompletionParams{
			Type: stripe.String("hosted_confirmation"),
			HostedConfirmation: &stripe.PaymentLinkCreateAfterCompletionHostedConfirmationParams{
				CustomMessage: stripe.String("🎉 Danke für das Aufladen deines Kontos. Du kannst den Browser nun schließen. Jeden moment sollte deine aufladung verarbeitet sein 💸"),
			},
		},
		LineItems: []*stripe.PaymentLinkCreateLineItemParams{
			&stripe.PaymentLinkCreateLineItemParams{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.PaymentLinkCreateLineItemPriceDataParams{
					Currency: stripe.String("EUR"),
					ProductData: &stripe.PaymentLinkCreateLineItemPriceDataProductDataParams{
						Name: stripe.String("Konto aufladung von 5 Euro"),
					},
					UnitAmount: stripe.Int64(500),
				},
			},
		},
	}

	result, err := s.StripeClient.V1PaymentLinks.Create(context.TODO(), params)
	if err != nil {
		log.Fatalf("error calling stripe: %v", err)
	}

	slog.Debug("stripe result", "url", result.URL, "result", result)

	return result.URL, nil

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
