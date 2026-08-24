package chargementservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/KochKevin/effective-spoon-v2/internal/chargements"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v86"
)

type BalanceChargementIntentCache interface {
	GetBalanceChargementIntentId() uuid.UUID
	ClearBalanceChargementIntentId()
}

const stripeEventCheckoutSessionCompleted = "checkout.session.completed"

func (s *Service) Ticker(ctx context.Context) {

	ticker := time.NewTicker(10 * time.Second)

	go func() {
		defer ticker.Stop()
		slog.Info("ticker started")

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				slog.Debug("ticker tick")
				err := s.stripeEventCheckoutCompletedChecker(ctx)
				if err != nil {
					slog.Error("stripe event checkout complete checker failed: %w", err)
				}
			}
		}

	}()

}

func (s *Service) stripeEventCheckoutCompletedChecker(ctx context.Context) error {

	lastEventId, err := s.Repo.GetStripeLastEventId(ctx)
	if err != nil {
		slog.Error("error getting the stripe last event id")

	}

	params := &stripe.EventListParams{Type: stripe.String(stripeEventCheckoutSessionCompleted)}
	if len(lastEventId) > 0 {
		params.EndingBefore = stripe.String(lastEventId)
	}
	params.Limit = stripe.Int64(1000)
	result := s.StripeClient.V1Events.List(ctx, params)

	events := result.Data()

	if len(events) > 0 {
		slog.Info("stripe event checkout timer processing", "count", len(events))
	}

	for _, event := range events {

		lastEventId = event.ID

		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			slog.Error("error while unmarshaling the "+stripeEventCheckoutSessionCompleted, "err", err, "event_id", event.ID)
			continue
		}

		if len(session.Metadata) == 0 {
			slog.Error("stripe event metadata is empty")
			continue
		}

		chargementIntentIdStr, ok := session.Metadata[MetadataBalanceChargementIntentId]
		if !ok {
			slog.Warn("no "+MetadataBalanceChargementIntentId+" in stripe event metadata "+stripeEventCheckoutSessionCompleted+" found, skipping event", "event_id", event.ID)
			continue
		}

		chargementIntentId, err := uuid.Parse(chargementIntentIdStr)
		if err != nil {
			slog.Error("invalid UUID in metadata", "err", err, "uuid_str_got", chargementIntentIdStr)
			continue
		}

		slog.Debug("chargement intent id", "id", chargementIntentId)

		var chargementIntent chargements.ChargementIntent

		err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

			chargementIntent, err = s.Repo.GetBalanceChargementIntent(ctx, tx, chargementIntentId)
			if err != nil {
				//slog.Error("invalid UUID in metadata", "err", err, "uuid_str_got", chargementIntentIdStr)
				return fmt.Errorf("invalid UUID in metadata: %w", err)
			}

			stripeCheckoutId := session.ID

			transaction, err := chargementIntent.Complete(stripeCheckoutId)
			if errors.Is(err, chargements.ErrAlreadyCompleted) {
				//slog.Error("chargement already completed")
				return err
			}

			if err != nil {
				slog.Error("error", err)
				return err
			}

			_, err = s.UserRepo.CreateTransaction(ctx, tx, transaction)
			if err != nil {
				//slog.Error("transaction can not be saved", chargementIntentIdStr)
				return fmt.Errorf("transaction obj can not be created: %w", err)
			}

			err = s.Repo.SaveBalanceChargementIntent(ctx, tx, chargementIntent)
			if err != nil {
				slog.Error("error", err)
				return err
			}

			err = s.Repo.SetStripeLastEventId(ctx, tx, lastEventId)
			if err != nil {
				//slog.Error("error setting the last stripe event id ")
				return fmt.Errorf("error setting the last stripe event id: %w", err)
			}

			// Set Payment link Invalid
			s.deactivateStripePaymentLink(ctx, session.PaymentLink.ID)

			return nil

		})

		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Empty sql result, mainly bevause those events werent created in the db", err)
			continue
		}

		if err != nil {
			slog.Error("error in transaction", err)
			return err
		}

		if chargementIntent.Id == s.BalanceChargementIntentCache.GetBalanceChargementIntentId() {
			//Notify that an chargment sucessfully happend and that the frontend can update
			s.BalanceChargementIntentCache.ClearBalanceChargementIntentId()
		}

	}

	return nil
}
