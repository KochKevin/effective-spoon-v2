package users

import (
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/google/uuid"
)

type Usercode string

type User struct {
	Id      uuid.UUID
	Name    string
	Balance money.Money
	Code    Usercode
}

func UserFrom(id uuid.UUID, name string, balance money.Money, usercode Usercode) User {
	return User{
		Id:      id,
		Name:    name,
		Balance: balance,
		Code:    usercode,
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
		Id:     id,
		UserID: userId,
		Amount: amount,
	}
}
