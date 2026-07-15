package users

import (
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/google/uuid"
)

type User struct {
	Id      uuid.UUID
	Name    string
	Balance money.Money
}

func UserFrom(id uuid.UUID, name string, balance money.Money) User {
	return User{
		Id:      id,
		Name:    name,
		Balance: balance,
	}
}

type Transaction struct {
	Id     uuid.UUID
	UserID uuid.UUID
	Amount money.Money
}

func NewTransaction(userID uuid.UUID, amount money.Money) Transaction {
	return Transaction{
		Id:     uuid.New(),
		UserID: userID,
		Amount: amount,
	}
}

func TranscationFrom(id uuid.UUID, userId uuid.UUID, amount money.Money) Transaction {
	return Transaction{
		Id: id,
		UserID: userId,
		Amount: amount,
	}
}