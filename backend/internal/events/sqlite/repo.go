package eventssqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
)

//r *Repo github.com/KochKevin/effective-spoon-v2/internal/events/Service.Repo

type Repo struct {
	Queries sqlc.Queries
}

func (r *Repo) CreateEvent(ctx context.Context, tx *sql.Tx, event events.Event) error {

	_, err := r.Queries.WithTx(tx).CreateEvent(ctx, sqlc.CreateEventParams{
		ID:                        event.Id,
		UserID:                    event.UserId,
		AmountFreeProductsPerUser: int64(event.AmountFreeProductsPerUser),
		StartTimestamp:            event.StartTimestamp,
		EndTimestamp:              event.EndTimestamp,
		Status:                    string(event.Status),
	})

	if err != nil {
		return fmt.Errorf("in creating event query: %w", err)
	}

	return nil

}
func (r *Repo) GetEventUsage(ctx context.Context, tx *sql.Tx, eventId uuid.UUID, userId uuid.UUID) (events.EventUsage, error) {

	eventUsage, err := r.Queries.WithTx(tx).GetEventUsage(ctx, sqlc.GetEventUsageParams{
		UserID:  userId,
		EventID: eventId,
	})

	// No event usage was created for the user, create it now
	if err == sql.ErrNoRows {
		return r.CreateEventUsage(ctx, tx, eventId, userId)
	}

	if err != nil {
		return events.EventUsage{}, fmt.Errorf("in getting event usage query: %w", err)
	}

	return events.EventUsage{
		UserId:     eventUsage.UserID,
		EventId:    eventUsage.EventID,
		AmountUsed: int(eventUsage.UsedAmountFreeProducts)}, nil
}

func (r *Repo) CreateEventUsage(ctx context.Context, tx *sql.Tx, eventId uuid.UUID, userId uuid.UUID) (events.EventUsage, error) {
	eventUsage, err := r.Queries.WithTx(tx).UpsertEventUsage(ctx, sqlc.UpsertEventUsageParams{
		UserID:                 userId,
		EventID:                eventId,
		UsedAmountFreeProducts: 0,
	})

	if err != nil {
		return events.EventUsage{}, fmt.Errorf("in creating event usage query: %w", err)
	}

	return events.EventUsage{
		UserId:     eventUsage.UserID,
		EventId:    eventUsage.EventID,
		AmountUsed: int(eventUsage.UsedAmountFreeProducts)}, nil
}

func (r *Repo) GetEvent(ctx context.Context, tx *sql.Tx, id uuid.UUID) (events.Event, error) {
	event, err := r.Queries.WithTx(tx).GetEvent(ctx, id)

	if err != nil {
		return events.Event{}, fmt.Errorf("in getting event query: %w", err)
	}

	return events.Event{
		Id:                        event.ID,
		UserId:                    event.UserID,
		AmountFreeProductsPerUser: int(event.AmountFreeProductsPerUser),
		StartTimestamp:            event.StartTimestamp,
		EndTimestamp:              event.EndTimestamp,
		Status:                    events.State(event.Status),
	}, nil
}

func (r *Repo) GetActiveEvent(ctx context.Context, tx *sql.Tx) (events.Event, error) {
	event, err := r.Queries.WithTx(tx).GetEventByStatus(ctx, string(events.Active))

	if err != nil {
		return events.Event{}, fmt.Errorf("in getting event query: %w", err)
	}

	return events.Event{
		Id:                        event.ID,
		UserId:                    event.UserID,
		AmountFreeProductsPerUser: int(event.AmountFreeProductsPerUser),
		StartTimestamp:            event.StartTimestamp,
		EndTimestamp:              event.EndTimestamp,
		Status:                    events.State(event.Status),
	}, nil
}

/*
func (r *Repo) SetStripeLastEventId(ctx context.Context, tx *sql.Tx, lastEventId string) error {

	err := r.Queries.WithTx(tx).SetLastStripeEventId(ctx, lastEventId)
	if err != nil {
		slog.Error("Error in set stripe last event id query", err)
		return err
	}

	return nil
}

func (r *Repo) GetStripeLastEventId(ctx context.Context) (string, error) {
	lastEventId, err := r.Queries.GetLastStripeEventId(ctx)
	if err != nil {
		slog.Error("Error in get stripe last event id query", err)
		return "", err
	}

	return lastEventId, nil
}

func (r *Repo) GetBalanceChargementIntent(ctx context.Context, tx *sql.Tx, id uuid.UUID) (chargements.ChargementIntent, error) {
	chargement, err := r.Queries.WithTx(tx).GetChargementIntent(ctx, id)

	if err != nil {
		slog.Error("Error in get chargement query", err)
		return chargements.ChargementIntent{}, err
	}

	return chargements.ChargementIntent{Id: chargement.ID, Status: chargements.Status(chargement.Status), User: chargement.UserID, Amount: money.MoneyFrom(int(chargement.Amount)), StripeCheckoutId: *chargement.StripeCheckoutID, TransactionId: chargement.TransactionID, PaymentLink: *chargement.StripePaymentLink, PaymentLinkId: *chargement.StripePaymentLinkID}, nil
}

func (r *Repo) SaveBalanceChargementIntent(ctx context.Context, tx *sql.Tx, chargementIntent chargements.ChargementIntent) error {
	err := r.Queries.WithTx(tx).UpdateChargementIntent(ctx, sqlc.UpdateChargementIntentParams{
		ID:                  chargementIntent.Id,
		Status:              string(chargementIntent.Status),
		UserID:              chargementIntent.User,
		Amount:              int64(chargementIntent.Amount.Cents),
		StripeCheckoutID:    &chargementIntent.StripeCheckoutId,
		TransactionID:       chargementIntent.TransactionId,
		StripePaymentLink:   &chargementIntent.PaymentLink,
		StripePaymentLinkID: &chargementIntent.PaymentLinkId,
	})
	if err != nil {
		slog.Error("Error in save chargement query", err)
		return err
	}

	return nil
}

func (r *Repo) CreateBalanceChargementIntent(ctx context.Context, tx *sql.Tx, chargementIntent chargements.ChargementIntent) error {

	_, err := r.Queries.WithTx(tx).CreateChargementIntent(ctx, sqlc.CreateChargementIntentParams{
		ID:                  chargementIntent.Id,
		Status:              string(chargementIntent.Status),
		UserID:              chargementIntent.User,
		Amount:              int64(chargementIntent.Amount.Cents),
		StripeCheckoutID:    &chargementIntent.StripeCheckoutId,
		TransactionID:       chargementIntent.TransactionId,
		StripePaymentLink:   &chargementIntent.PaymentLink,
		StripePaymentLinkID: &chargementIntent.PaymentLinkId,
	})
	if err != nil {
		slog.Error("Error in create chargement query", err)
		return err
	}

	return nil

}
*/
