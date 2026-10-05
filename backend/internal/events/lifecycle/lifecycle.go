package eventslifecycle

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
)

//Start Event
// - Load Event
//End Event

type EventService interface {
	CreateCurrentEventTx(ctx context.Context, tx *sql.Tx, authorId uuid.UUID, amountPerPerson int, endDateTime time.Time) (events.Event, error)
}

type ShoppingCartService interface {
	AssignEventToCurrentShoppingCartTx(ctx context.Context, tx *sql.Tx, userId uuid.UUID, eventId uuid.UUID) error
}

type Service struct {
	Txm                 infrastructure.TxManager
	EventService        EventService
	ShoppingCartService ShoppingCartService
}

func NewService(txm infrastructure.TxManager, eventService EventService, shoppingCartService ShoppingCartService) *Service {

	return &Service{
		Txm:                 txm,
		EventService:        eventService,
		ShoppingCartService: shoppingCartService,
	}
}

// Check if an event us current active and load it into cache
/*
func (s *Service) Initalize() {

}
*/

// Create new Event, set it as the current active Event and schedule it ending
func (s *Service) StartCurrentEvent(ctx context.Context, authorId uuid.UUID, amountPerPerson int, endDateTime time.Time) (event events.Event, err error) {

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		event, err = s.EventService.CreateCurrentEventTx(ctx, tx, authorId, amountPerPerson, endDateTime)
		if err != nil {
			return fmt.Errorf("creating current event: %w", err)
		}

		err = s.ShoppingCartService.AssignEventToCurrentShoppingCartTx(ctx, tx, authorId, event.Id)
		if err != nil {
			return fmt.Errorf("assigning event to current shopping cart: %w", err)
		}

		return nil
	})

	if err != nil {
		return events.Event{}, fmt.Errorf("in transaction: %w", err)
	}

	time.AfterFunc(time.Until(event.EndTimestamp), func() {
		s.End(event)
	})

	return event, nil
}

func (s *Service) End(event events.Event) {

}
