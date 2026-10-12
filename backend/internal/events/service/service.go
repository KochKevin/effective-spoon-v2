package eventsservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
)

type Repo interface {
	CreateEvent(ctx context.Context, tx *sql.Tx, event events.Event) error
	GetEvent(ctx context.Context, tx *sql.Tx, id uuid.UUID) (events.Event, error)
	GetActiveEvent(ctx context.Context, tx *sql.Tx) (events.Event, error)
	SaveEvent(ctx context.Context, tx *sql.Tx, event events.Event) error

	// Query all Event Usages of this event, sum the used amounts per product and calculate the resulting total cost
	GetTotalCostsOfEvent(ctx context.Context, tx *sql.Tx, id uuid.UUID) (money.Money, error)
}

type CacheRepo interface {
	SetCurrentEventId(id uuid.UUID)
	GetCurrentEventId() uuid.UUID
	ClearCurrentEventId()
}

type UserService interface {
	//TODO: Use this Service function in Shoppingcart too
	CreateTransactionTx(ctx context.Context, tx *sql.Tx, transaction users.Transaction) error
}

type Service struct {
	Txm         infrastructure.TxManager
	Repo        Repo
	CacheRepo   CacheRepo
	UserService UserService
}

// Setup timer to load current event and end it after end time

func NewService(txm infrastructure.TxManager, repo Repo, cacheRepo CacheRepo, userService UserService) Service {

	return Service{
		Txm:         txm,
		Repo:        repo,
		CacheRepo:   cacheRepo,
		UserService: userService,
	}

}
func (s *Service) LoadCurrentlyActiveEventIntoCache(ctx context.Context) (event events.Event, err error) {
	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		event, err = s.Repo.GetActiveEvent(ctx, tx)

		if err != nil {
			return fmt.Errorf("error getting currently active event from persitent volume: %w", err)
		}

		//Set Current
		go s.CacheRepo.SetCurrentEventId(event.Id)

		return nil
	})
	if err != nil {
		return events.Event{}, fmt.Errorf("in transaction: %w", err)
	}

	return event, nil
}

// Only one event at the time is allowed to exist
func (s *Service) CreateCurrentEvent(ctx context.Context, authorId uuid.UUID, amountPerPerson int, endDateTime time.Time) (event events.Event, err error) {

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		event, err = s.CreateCurrentEventTx(ctx, tx, authorId, amountPerPerson, endDateTime)
		if err != nil {
			return fmt.Errorf("in create current event tx: %w", err)
		}

		return nil
	})
	if err != nil {
		return events.Event{}, fmt.Errorf("in transaction: %w", err)
	}

	return event, nil

}

func (s *Service) CreateCurrentEventTx(ctx context.Context, tx *sql.Tx, authorId uuid.UUID, amountPerPerson int, endDateTime time.Time) (events.Event, error) {

	//Only proceed if no other event is currently loaded
	if uuid.Nil != s.CacheRepo.GetCurrentEventId() {
		return events.Event{}, fmt.Errorf("there is already an event set in cache. Only one event can exist at one time")
	}

	event, err := events.NewEvent(authorId, amountPerPerson, time.Now(), endDateTime, events.Active)
	if err != nil {
		return events.Event{}, err
	}

	err = s.Repo.CreateEvent(ctx, tx, event)

	if err != nil {
		return events.Event{}, fmt.Errorf("error creating event on persitent volume: %w", err)
	}

	//Set Current
	s.CacheRepo.SetCurrentEventId(event.Id)

	return event, nil
}

func (s *Service) GetCurrentEvent(ctx context.Context) (event events.Event, err error) {

	//Only proceed if an event is currently loaded
	if uuid.Nil == s.CacheRepo.GetCurrentEventId() {
		return events.Event{}, events.ErrNoCurrentEvent
	}

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		event, err = s.Repo.GetEvent(ctx, tx, s.CacheRepo.GetCurrentEventId())

		if err != nil {
			return fmt.Errorf("error getting event from persitent volume: %w", err)
		}

		return nil
	})
	if err != nil {
		return events.Event{}, fmt.Errorf("in transaction: %w", err)
	}

	return event, nil

}

// Remove need for transaction - GET Repos dont need an transaction
func (s *Service) GetEvent(ctx context.Context, eventId uuid.UUID) (event events.Event, err error) {

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		event, err = s.Repo.GetEvent(ctx, tx, eventId)

		if err != nil {
			return fmt.Errorf("error getting event from persitent volume: %w", err)
		}

		return nil
	})
	if err != nil {
		return events.Event{}, fmt.Errorf("in transaction: %w", err)
	}

	return event, nil
}

// Invalidate Event for further usage, generate Transaction and clear current event cache
func (s *Service) Checkout(ctx context.Context, eventId uuid.UUID) (err error) {
	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		event, err := s.Repo.GetEvent(ctx, tx, eventId)
		if err != nil {
			return fmt.Errorf("in get event: %w", err)
		}

		err = event.Checkout()
		if err != nil {
			return fmt.Errorf("in checkout: %w", err)
		}

		err = s.Repo.SaveEvent(ctx, tx, event)
		if err != nil {
			return fmt.Errorf("in save event: %w", err)
		}

		err = s.createTransactionFromEvent(ctx, tx, event)
		if err != nil {
			return fmt.Errorf("in create transaction from event: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("in transaction: %w", err)
	}

	s.CacheRepo.ClearCurrentEventId()

	return nil
}

func (s *Service) createTransactionFromEvent(ctx context.Context, tx *sql.Tx, event events.Event) error {

	cost, err := s.Repo.GetTotalCostsOfEvent(ctx, tx, event.Id)
	if err != nil {
		return fmt.Errorf("in get total cost of event: %w", err)
	}

	err = s.UserService.CreateTransactionTx(ctx, tx, users.NewTransaction(event.UserId, cost))
	if err != nil {
		return fmt.Errorf("in creating user.transaction: %w", err)
	}

	return nil
}
