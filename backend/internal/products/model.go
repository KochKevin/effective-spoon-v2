package products

import "github.com/google/uuid"

type Product struct{
	Id uuid.UUID
	Name string
	Price int
}

func NewProduct(id uuid.UUID, name string, price int) Product {
	return Product{
		Id: id,
		Name: name,
		Price: price,
	}
}