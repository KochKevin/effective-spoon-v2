package authcache

import (
	"sync"

	"github.com/google/uuid"
)

type Cache struct {
	mutex  sync.RWMutex
	userId uuid.UUID
}

func New() *Cache {
	return &Cache{
		userId: uuid.Nil,
	}
}

func (c *Cache) SetCurrentUserId(userId uuid.UUID) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.userId = userId
}

func (c *Cache) GetCurrentUserId() (userId uuid.UUID) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.userId
}

func (c *Cache) TryLogin(userId uuid.UUID) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.userId == uuid.Nil {
		c.userId = userId
		return true
	}

	return false
}
