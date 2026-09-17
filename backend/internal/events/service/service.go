package eventsservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
)

type Repo interface {
	CreateEvent(ctx context.Context, tx *sql.Tx, event events.Event) error
	GetEvent(ctx context.Context, tx *sql.Tx, id uuid.UUID) (events.Event, error)
	GetEventUsage(ctx context.Context, tx *sql.Tx, eventId uuid.UUID, userId uuid.UUID) (events.EventUsage, error)
}

type CacheRepo interface {
	SetCurrentEventId(id uuid.UUID)
	GetCurrentEventId() uuid.UUID
	ClearCurrentEventId()
}

type Service struct {
	Txm       infrastructure.TxManager
	Repo      Repo
	CacheRepo CacheRepo
}

// Setup timer to load current event and end it after end time

func NewService(txm infrastructure.TxManager, repo Repo, cacheRepo CacheRepo) (Service, error) {

	return Service{
		Txm:       txm,
		Repo:      repo,
		CacheRepo: cacheRepo,
	}, nil

}

// Only one event at the time is allowed to exist
func (s *Service) CreateCurrentEvent(ctx context.Context, authorId uuid.UUID, amountPerPerson int, endDateTime time.Time) (events.Event, error) {

	//Only proceed if no other event is currently loaded
	if uuid.Nil != s.CacheRepo.GetCurrentEventId() {
		return events.Event{}, fmt.Errorf("there is already an event set in cache. Only one event can exist at one time: %w")
	}

	event, err := events.NewEvent(authorId, amountPerPerson, time.Now(), endDateTime, events.Active)
	if err != nil {
		return events.Event{}, err
	}

	err = s.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		err := s.Repo.CreateEvent(context.Background(), tx, event)

		if err != nil {
			return fmt.Errorf("error creating event on persitent volume: %w", err)
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

func (s *Service) GetEventUsage(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (eventUsage events.EventUsage, err error) {

	err = s.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		eventUsage, err = s.Repo.GetEventUsage(context.Background(), tx, eventId, userId)

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

func (s *Service) GetCurrentEvent(ctx context.Context) (event events.Event, err error) {

	//Only proceed if an event is currently loaded
	if uuid.Nil == s.CacheRepo.GetCurrentEventId() {
		return events.Event{}, events.NoCurrentEventErr
	}

	err = s.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

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

func (s *Service) endEvent() {}
