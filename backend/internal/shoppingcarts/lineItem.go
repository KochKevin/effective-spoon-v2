package shoppingcarts

import (
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
)

type LineItem struct {
	Product products.Product
	Amount  Amount
}

func (l *LineItem) GetPrice() money.Money {
	//slog.Debug("GetPrice of lineItem", "price", money.MoneyFrom(l.Product.Price * l.Amount))
	return l.Product.Price.Multi(l.Amount.GetPayedAmount())
}

/*
	func NewLineItem(product products.Product) LineItem {
		return LineItem{
			Product: product,
			Amount:  NewAmount(),
		}
	}
*/
