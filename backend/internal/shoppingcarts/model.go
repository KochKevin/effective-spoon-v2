package shoppingcarts

import (
	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/google/uuid"
)

type ShoppingCart struct {
	Id        uuid.UUID
	LineItems []LineItem
}

func (s *ShoppingCart) GetFullPrice() money.Money {
	total := money.MoneyFrom(0)

	//slog.Debug("GetFullPrice of shopping cart ", "total", total)

	for _, item := range s.LineItems {
		total = total.Add(item.GetPrice())
		//slog.Debug("GetFullPrice", "total", total)
	}
	//slog.Debug("Final GetGullPrice", "total", total)
	return total
}

// Increase the amount of an line item, if no product is found. It will be added
func (s *ShoppingCart) IncreaseProductAmount(product products.Product) {

	slog.Debug("Searching for product to increase: ", product.Id)

	for i := range s.LineItems {

		if s.LineItems[i].Product.Id == product.Id {
			slog.Debug("Increase amount on product", product.Id.String())
			s.LineItems[i].IncreaseAmount()
			return
		}

	}

	//slog.Error("Did not found product ", productId.String(), " in shopping cart ", s.Id.String(), " to increase it")
	s.AddProduct(product)
}

// Add completly new product to line items
func (s *ShoppingCart) AddProduct(product products.Product) {
	s.LineItems = append(s.LineItems, NewLinteItem(product, 1))
}

func ShoppingCartFrom(id uuid.UUID, lineItems []LineItem) ShoppingCart {
	return ShoppingCart{
		Id:        id,
		LineItems: lineItems,
	}
}

func NewShoppingCart() ShoppingCart {
	return ShoppingCart{
		Id:        uuid.New(),
		LineItems: nil,
	}
}

type LineItem struct {
	Product products.Product
	Amount  int
}

func (l *LineItem) GetPrice() money.Money {
	//slog.Debug("GetPrice of lineItem", "price", money.MoneyFrom(l.Product.Price * l.Amount))
	return money.MoneyFrom(l.Product.Price * l.Amount)
}

func (l *LineItem) IncreaseAmount() {
	l.Amount++
}

func NewLinteItem(product products.Product, amount int) LineItem {
	return LineItem{
		Product: product,
		Amount:  amount,
	}
}
