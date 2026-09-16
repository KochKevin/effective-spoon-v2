package userservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
)

type Repo interface {
	GetUser(ctx context.Context, tx *sql.Tx, id uuid.UUID) (users.User, error)
}

type Service struct {
	Txm  infrastructure.TxManager
	Repo Repo
}

func (s *Service) GetUser(ctx context.Context, userId uuid.UUID) (user users.User, err error) {

	err = s.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		user, err = s.Repo.GetUser(ctx, tx, userId)

		if err != nil {
			return fmt.Errorf("getting user: %w", err)
		}

		return nil

	})

	if err != nil {
		return users.User{}, fmt.Errorf("in transaction: %w", err)
	}

	return user, nil

}
