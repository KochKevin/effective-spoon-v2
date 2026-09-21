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
	ShoppingCartCanceled   ShoppingCartStatus = "canceled"
	ShoppingCartCheckedOut ShoppingCartStatus = "checked-out"
)

type ShoppingCart struct {
	Id            uuid.UUID
	LineItems     []LineItem
	UserId        uuid.UUID
	TransactionId uuid.NullUUID
	Status        ShoppingCartStatus

	//Shopping Carts can be connected to an ongoing event to use an contigent of free products
	UseEvent bool
	EventId  uuid.UUID
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
func (s *ShoppingCart) IncreaseProductAmount(product products.Product, addAsFreeProduct bool) {

	slog.Debug("Searching for product to increase: ", product.Id)

	for i := range s.LineItems {

		if s.LineItems[i].Product.Id == product.Id {
			slog.Debug("Increase amount on product", product.Id.String())

			if addAsFreeProduct {
				s.LineItems[i].Amount.IncreaseFreeAmount()
			} else {
				s.LineItems[i].Amount.IncreasePayedAmount()
			}
			return
		}

	}

	//slog.Error("Did not found product ", productId.String(), " in shopping cart ", s.Id.String(), " to increase it")
	//Fallback if the product to increase is not found in the shopping cart
	if addAsFreeProduct {
		s.LineItems = append(s.LineItems, LineItem{
			Product: product,
			Amount:  AmountFrom(0, 1),
		})
	} else {
		s.LineItems = append(s.LineItems, LineItem{
			Product: product,
			Amount:  AmountFrom(1, 0),
		})
	}
}
func (s *ShoppingCart) DecreaseProductAmount(productId uuid.UUID) {

	slog.Debug("Searching for product to decrease: ", productId)

	for i := range s.LineItems {

		if s.LineItems[i].Product.Id == productId {

			// Remove LineItems when there is only one left
			if s.LineItems[i].Amount.GetTotalAmount() <= 1 {
				slog.Debug("Found product to decrese, but it needs to be removed")
				s.LineItems = append(s.LineItems[:i], s.LineItems[i+1:]...)
			} else {
				slog.Debug("Found product to decrese, decresing it")
				s.LineItems[i].Amount.DecreaseAmount()
			}
			return

		}

	}

	slog.Error("Did not found product ", productId.String(), " in shopping cart ", s.Id.String(), " to decrease it")

}

//TODO: Combine Generate Transaction and Checkout to one Checkout Function which returns an transaction

// Creates a new Transaction for the buying user based on the data in the shoppingcart
func (s *ShoppingCart) GenerateTransaction() (transaction users.Transaction) {

	slog.Debug("Generate Transaction", "shopping cart Full price", s.GetFullPrice())

	slog.Debug("Generate Transaction", "transaction amount", money.MoneyFrom(0).Sub(s.GetFullPrice()))

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

/*
// Calculate the amount of free products used in this shopping cart
func (s *ShoppingCart) GetAmountFreeProductsUsed() int {
	count := 0

	for _, item := range s.LineItems {
		count += item.Amount.GetFreeAmount()
	}

	return count
}
*/

func NewShoppingCart(userId uuid.UUID, eventId uuid.UUID, useEvent bool) ShoppingCart {
	return ShoppingCart{
		Id:            uuid.New(),
		LineItems:     nil,
		UserId:        userId,
		TransactionId: uuid.NullUUID{Valid: false},
		Status:        ShoppingCartActive,
		EventId:       eventId,
		UseEvent:      useEvent,
	}
}

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
type Amount struct {
	payedAmount int
	freeAmount  int
}

/*
	func NewAmount() Amount {
		return Amount{
			payedAmount: 0,
			freeAmount:  0,
		}
	}
*/
func AmountFrom(payedAmount int, freeAmount int) Amount {
	return Amount{
		payedAmount: payedAmount,
		freeAmount:  freeAmount,
	}
}

func (a *Amount) GetFreeAmount() int {
	return a.freeAmount
}

func (a *Amount) GetPayedAmount() int {
	return a.payedAmount
}

func (a *Amount) GetTotalAmount() int {
	return a.payedAmount + a.freeAmount
}

func (a *Amount) IncreaseFreeAmount() {
	a.freeAmount++

}

func (a *Amount) IncreasePayedAmount() {
	a.payedAmount++
}

// First decrease payed products then free products
func (a *Amount) DecreaseAmount() {

	if a.payedAmount > 0 {
		a.payedAmount--
	} else if a.freeAmount > 0 {
		a.freeAmount--
	}
}
