package chargementssqlite

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/chargements"
	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/google/uuid"
)

//r *Repo github.com/KochKevin/effective-spoon-v2/internal/chargements/Service.Repo

type Repo struct {
	Queries sqlc.Queries
}

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

	return chargements.ChargementIntent{Id: chargement.ID, Status: chargements.Status(chargement.Status), User: chargement.UserID, Amount: money.MoneyFrom(int(chargement.Amount)), StripeCheckoutId: *chargement.StripeCheckoutID, TransactionId: chargement.TransactionID}, nil
}

func (r *Repo) SaveBalanceChargementIntent(ctx context.Context, tx *sql.Tx, chargementIntent chargements.ChargementIntent) error {
	err := r.Queries.WithTx(tx).UpdateChargementIntent(ctx, sqlc.UpdateChargementIntentParams{
		ID:               chargementIntent.Id,
		Status:           string(chargementIntent.Status),
		UserID:           chargementIntent.User,
		Amount:           int64(chargementIntent.Amount.Cents),
		StripeCheckoutID: &chargementIntent.StripeCheckoutId,
		TransactionID:    chargementIntent.TransactionId,
	})
	if err != nil {
		slog.Error("Error in save chargement query", err)
		return err
	}

	return nil
}

func (r *Repo) CreateBalanceChargementIntent(ctx context.Context, tx *sql.Tx, chargementIntent chargements.ChargementIntent) error {

	_, err := r.Queries.WithTx(tx).CreateChargementIntent(ctx, sqlc.CreateChargementIntentParams{
		ID:               chargementIntent.Id,
		Status:           string(chargementIntent.Status),
		UserID:           chargementIntent.User,
		Amount:           int64(chargementIntent.Amount.Cents),
		StripeCheckoutID: &chargementIntent.StripeCheckoutId,
		TransactionID:    chargementIntent.TransactionId,
	})
	if err != nil {
		slog.Error("Error in create chargement query", err)
		return err
	}

	return nil

}
