package products

import (
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/google/uuid"
)

type Product struct{
	Id uuid.UUID
	Name string
	Price money.Money
}

func NewProduct(id uuid.UUID, name string, price money.Money) Product {
	return Product{
		Id: id,
		Name: name,
		Price: price,
	}
}