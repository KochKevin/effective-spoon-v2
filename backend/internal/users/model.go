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
