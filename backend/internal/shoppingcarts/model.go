package shoppingcarts

import (
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/google/uuid"
)

type ShoppingCart struct {
	Id        uuid.UUID
	lineItems []lineItem
}

func (s *ShoppingCart) GetFullPrice() money.Money {
	var total money.Money


	for _, item := range s.lineItems{
		total.Add(item.GetPrice())
	}

	return total
}

func ShoppingCartFrom(id uuid.UUID, lineItems []lineItem) ShoppingCart {
	return ShoppingCart{
		Id:        id,
		lineItems: lineItems,
	}
}

func NewShoppingCart() ShoppingCart {
	return ShoppingCart{
		Id:        uuid.New(),
		lineItems: nil,
	}
}

type lineItem struct {
	Product products.Product
	Amount  int
}

func (l *lineItem) GetPrice() money.Money {
	return money.MoneyFrom(l.Product.Price * l.Amount)
}

func NewLinteItem(product products.Product, amount int) lineItem {
	return lineItem{
		Product: product,
		Amount:  amount,
	}
}
