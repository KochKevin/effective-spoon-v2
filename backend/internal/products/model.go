package products

import (
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/google/uuid"
)

type ProductCode string

type Product struct {
	Id    uuid.UUID
	Name  string
	Price money.Money
	Code  ProductCode
}

func NewProduct(id uuid.UUID, name string, price money.Money, code ProductCode) Product {
	return Product{
		Id:    id,
		Name:  name,
		Price: price,
		Code:  code,
	}
}
