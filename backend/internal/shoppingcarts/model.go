package shoppingcarts

import (
	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/google/uuid"
)

type ShoppingCartStatus string

var (
	ShoppingCartActive     ShoppingCartStatus = "active"
	ShoppingCartCheckedOut ShoppingCartStatus = "checked-out"
)

type ShoppingCart struct {
	Id            uuid.UUID
	LineItems     []LineItem
	UserId        uuid.UUID
	TransactionId uuid.NullUUID
	Status        ShoppingCartStatus
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

func (s *ShoppingCart) DecreaseProductAmount(productId uuid.UUID) {

	slog.Debug("Searching for product to decrease: ", productId)

	for i := range s.LineItems {

		if s.LineItems[i].Product.Id == productId {

			//Remove LineItems when there is only one left
			if s.LineItems[i].Amount <= 1 {
				slog.Debug("Found product to decrese, but it needs to be removed")
				s.LineItems = append(s.LineItems[:i], s.LineItems[i+1:]...)
			} else {
				slog.Debug("Found product to decrese, decresing it")
				s.LineItems[i].DecreaseAmount()
			}
			return

		}

	}

	slog.Error("Did not found product ", productId.String(), " in shopping cart ", s.Id.String(), " to decrease it")

}

// Creates a new Transaction for the buying user based on the data in the shoppingcart
func (s *ShoppingCart) GenerateTransaction() (transaction users.Transaction) {
	return users.NewTransaction(
		s.UserId,
		money.MoneyFrom(0).Sub(s.GetFullPrice()), //Withdrawl money from user with this Transaction
	)
}

func (s *ShoppingCart) Checkout(transactionId uuid.UUID) {
	s.TransactionId = uuid.NullUUID{
		Valid: true,
		UUID:  transactionId,
	}

	s.Status = ShoppingCartCheckedOut
}

func ShoppingCartFrom(id uuid.UUID, lineItems []LineItem, userID uuid.UUID, transactionId uuid.NullUUID, status ShoppingCartStatus) ShoppingCart {
	return ShoppingCart{
		Id:            id,
		LineItems:     lineItems,
		UserId:        userID,
		TransactionId: transactionId,
		Status:        status,
	}
}

func NewShoppingCart(userId uuid.UUID) ShoppingCart {
	return ShoppingCart{
		Id:            uuid.New(),
		LineItems:     nil,
		UserId:        userId,
		TransactionId: uuid.NullUUID{Valid: false},
		Status:        ShoppingCartActive,
	}
}

type LineItem struct {
	Product products.Product
	Amount  int
}

func (l *LineItem) GetPrice() money.Money {
	//slog.Debug("GetPrice of lineItem", "price", money.MoneyFrom(l.Product.Price * l.Amount))
	return l.Product.Price.Multi(l.Amount)
}

func (l *LineItem) IncreaseAmount() {
	l.Amount++
}

func (l *LineItem) DecreaseAmount() {
	l.Amount--
}

func NewLinteItem(product products.Product, amount int) LineItem {
	return LineItem{
		Product: product,
		Amount:  amount,
	}
}
