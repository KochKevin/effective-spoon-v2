package events

import (
	"time"

	"github.com/google/uuid"
)

type State string

const (
	Active   State = "ACTIVE"
	Inactive State = "INACTIVE"
)

func NewEvent(userId uuid.UUID, amountFreeProductsPerUser int, startTimestamp time.Time, endTimestamp time.Time, status State) (Event, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return Event{}, err
	}

	return Event{
		Id:                        id,
		UserId:                    userId,
		AmountFreeProductsPerUser: amountFreeProductsPerUser,
		StartTimestamp:            startTimestamp,
		EndTimestamp:              endTimestamp,
		Status:                    status,
	}, nil
}

type Event struct {
	Id                        uuid.UUID
	UserId                    uuid.UUID
	AmountFreeProductsPerUser int
	StartTimestamp            time.Time
	EndTimestamp              time.Time
	Status                    State
}

type EventUsage struct {
	UserId     uuid.UUID
	EventId    uuid.UUID
	AmountUsed int
}
