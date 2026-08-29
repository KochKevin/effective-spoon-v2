package chargements

import (
	"errors"

	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusCanceld   = "canceld"
)

type ChargementIntent struct {
	Id               uuid.UUID
	Status           Status
	User             uuid.UUID
	Amount           money.Money
	StripeCheckoutId string
	TransactionId    uuid.UUID
	PaymentLink      string
	PaymentLinkId    string
}

func NewChargementIntent(amount money.Money, user uuid.UUID) (ChargementIntent, error) {

	id := uuid.New()

	return ChargementIntent{
		Id:               id,
		Status:           StatusPending,
		Amount:           amount,
		StripeCheckoutId: "",
		TransactionId:    uuid.Nil,
		User:             user,
		PaymentLink:      "empty",
		PaymentLinkId:    "empty",
	}, nil
}

func (c *ChargementIntent) AddStripePaymentLinkAndId(paymentLink string, paymentLinkId string) {
	c.PaymentLink = paymentLink
	c.PaymentLinkId = paymentLinkId
}

var ErrAlreadyCompleted = errors.New("this chargment intent is already completed")

// Complete cmpletes the chargment intent and creates and transaction for the user
func (c *ChargementIntent) Complete(stripeCheckoutId string) (users.Transaction, error) {

	if c.Status == StatusCompleted {
		return users.Transaction{}, ErrAlreadyCompleted
	}

	c.Status = StatusCompleted
	c.StripeCheckoutId = stripeCheckoutId

	transaction := users.NewTransaction(c.User, c.Amount)

	c.TransactionId = transaction.Id
	return transaction, nil
}
