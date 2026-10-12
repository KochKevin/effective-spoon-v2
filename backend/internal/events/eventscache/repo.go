package eventscache

import (
	"sync"

	"github.com/google/uuid"
)

type Cache struct {
	mutex          sync.RWMutex
	currentEventId uuid.UUID
}

func New() *Cache {
	return &Cache{
		currentEventId: uuid.Nil,
	}
}

func (c *Cache) SetCurrentEventId(id uuid.UUID) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.currentEventId = id
}

func (c *Cache) GetCurrentEventId() uuid.UUID {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.currentEventId
}

func (c *Cache) ClearCurrentEventId() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.currentEventId = uuid.Nil
}
