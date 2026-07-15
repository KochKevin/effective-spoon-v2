package userssqlite

import (
	"context"
	"database/sql"

	"log/slog"

	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/google/uuid"
)

type Repo struct {
	Queries sqlc.Queries
}

// TODO: calling Create should also create the line items
func (r *Repo) GetUser(ctx context.Context, tx *sql.Tx, id uuid.UUID) (users.User, error) {
	user, err := r.Queries.WithTx(tx).GetUser(ctx, id)
	if err != nil {
		slog.Error("Error in get user query", "error", err)
		return users.User{}, err
	}

	return users.UserFrom(user.ID, user.Name, money.MoneyFrom(int(user.Balance))), nil
}


func (r *Repo) CreateTransaction(ctx context.Context, tx *sql.Tx, transaction users.Transaction) (users.Transaction, error) {

	dbObj, err := r.Queries.WithTx(tx).CreateTransaction(ctx, sqlc.CreateTransactionParams{
		ID: transaction.Id,
		UserID: transaction.UserID,
		Amount: int64(transaction.Amount.GetAsCents()),
	})
	
	if err != nil {
		slog.Error("Error in get user query", "error", err)
		return users.Transaction{}, err
	}

	return users.TranscationFrom(dbObj.ID, dbObj.UserID, money.MoneyFrom(int(dbObj.Amount))), nil

}
