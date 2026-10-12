package eventslifecycle

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
)

type EventService interface {
	CreateCurrentEventTx(ctx context.Context, tx *sql.Tx, authorId uuid.UUID, amountPerPerson int, endDateTime time.Time) (events.Event, error)
	LoadCurrentlyActiveEventIntoCache(ctx context.Context) (event events.Event, err error)
	Checkout(ctx context.Context, eventId uuid.UUID) (err error)
}

type ShoppingCartService interface {
	AssignEventToCurrentShoppingCartTx(ctx context.Context, tx *sql.Tx, userId uuid.UUID, eventId uuid.UUID) error
}

type PushService interface {
	PushEventEnd()
}

type Service struct {
	Txm                 infrastructure.TxManager
	EventService        EventService
	ShoppingCartService ShoppingCartService
	PushService         PushService
}

func NewService(txm infrastructure.TxManager, eventService EventService, shoppingCartService ShoppingCartService, pushService PushService) Service {

	return Service{
		Txm:                 txm,
		EventService:        eventService,
		ShoppingCartService: shoppingCartService,
		PushService:         pushService,
	}
}

// Check if an event is should be currently active and load it from the DB into cache
func (s *Service) Initalize(ctx context.Context) error {
	event, err := s.EventService.LoadCurrentlyActiveEventIntoCache(ctx)
	if err != nil {
		return fmt.Errorf("loading currently active event into cache: %w", err)
	}

	s.ScheduleEnd(event.EndTimestamp, event.Id)

	return nil
}

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

		s.ScheduleEnd(event.EndTimestamp, event.Id)

		return nil
	})

	if err != nil {
		return events.Event{}, fmt.Errorf("in transaction: %w", err)
	}

	return event, nil
}

// Synctest this!
// at - when should the event end? eventId - which event
func (s *Service) ScheduleEnd(at time.Time, eventID uuid.UUID) {

	slog.Info("scheduling event end",
		"event_id", eventID,
		"at", at.UTC().String(),
		"now", time.Now().UTC().String(),
		"until", time.Until(at).String(),
	)

	time.AfterFunc(time.Until(at), func() {
		err := s.End(eventID)
		if err != nil {
			slog.Error("error in ending scheduled event", "err", err)
		}
	})

}

func (s *Service) End(eventId uuid.UUID) error {
	/*
		- Event State ändern zu inaktiv oder ended
			- Es sollte dann nicht mehr möglich sein das Event zu nutzen
		- Transaktion erstellen
			- Transaktion für alle Produkte erstellen und dem erstellten Nutzer in Rechnung stellen
		- SSE nachricht publishen, dass das Event vorbei ist
			- Frontend sollte seine stores leeren (Current Event und Event Usage)
	*/

	ctx := context.Background()

	err := s.EventService.Checkout(ctx, eventId)
	if err != nil {
		return fmt.Errorf("in event checkout: %w", err)
	}

	//Publish SSE Push

	s.PushService.PushEventEnd()

	return nil

}
