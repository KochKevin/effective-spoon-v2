package chargementcache

import (
	"sync"

	"github.com/google/uuid"
)

type Cache struct {
	mutex                     sync.RWMutex
	balanceChargementIntentId uuid.UUID
}

func New() *Cache {
	return &Cache{
		balanceChargementIntentId: uuid.Nil,
	}
}

func (c *Cache) SetBalanceChargementIntentId(id uuid.UUID) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.balanceChargementIntentId = id
}

func (c *Cache) GetBalanceChargementIntentId() uuid.UUID {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.balanceChargementIntentId
}

func (c *Cache) ClearBalanceChargementIntentId() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.balanceChargementIntentId = uuid.Nil
}
