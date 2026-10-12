package shoppingcarts

import (
	"github.com/KochKevin/effective-spoon-v2/internal/events"
)

type ShoppingCartView struct {
	Cart                 ShoppingCart
	EventUsage           events.EventUsage
	FreeProductPerPerson int
}
