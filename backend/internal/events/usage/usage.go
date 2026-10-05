package eventsusage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
)

type Repo interface {
	GetEventUsage(ctx context.Context, tx *sql.Tx, eventId uuid.UUID, userId uuid.UUID) (events.EventUsage, error)
}

type Usage struct {
	Txm  infrastructure.TxManager
	Repo Repo
}

func NewService(txm infrastructure.TxManager, repo Repo) Usage {

	service := Usage{
		Txm:  txm,
		Repo: repo,
	}

	return service

}

func (u *Usage) GetEventUsage(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (eventUsage events.EventUsage, err error) {

	err = u.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		eventUsage, err = u.Repo.GetEventUsage(ctx, tx, eventId, userId)

		if err != nil {
			return fmt.Errorf("getting event usage: %w", err)
		}

		return nil
	})
	if err != nil {
		return events.EventUsage{}, fmt.Errorf("in transaction: %w", err)
	}

	return eventUsage, nil
}
