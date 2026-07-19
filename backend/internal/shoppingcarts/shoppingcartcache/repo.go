package shoppingcartcache

import (
	"sync"

	"github.com/google/uuid"
)

type Cache struct {
	mutex          sync.RWMutex
	shoppingCartId uuid.UUID
}

func New() *Cache {
	return &Cache{
		shoppingCartId: uuid.Nil,
	}
}

func (c *Cache) SetCurrentCartId(cartId uuid.UUID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.shoppingCartId = cartId
}

func (c *Cache) GetCurrentCartId() (cartId uuid.UUID) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.shoppingCartId
}

func (c *Cache) ClearCurrentCartId() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.shoppingCartId = uuid.Nil
}
