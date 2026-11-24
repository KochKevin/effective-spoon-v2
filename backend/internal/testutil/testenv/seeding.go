package testenv

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/google/uuid"
)

// Create user and return it
func (e *Env) SeedUser(t *testing.T) (users.User, error) {

	user := users.NewUser("Test User", money.MoneyFrom(0), "usercode1")

	return user, e.App.Txm.WithTx(t.Context(), func(tx *sql.Tx) error {
		err := e.App.Repos.Users.CreateUser(
			t.Context(),
			tx,
			user,
		)

		if err != nil {
			return fmt.Errorf("in create user repo: %w", err)
		}

		return nil

	})

}

func (e *Env) SeedProduct(t *testing.T) (products.Product, error) {
	p := products.NewProduct(uuid.New(), "Test Product", money.MoneyFrom(100), "productcode1")

	return p, e.App.Txm.WithTx(t.Context(), func(tx *sql.Tx) error {
		err := e.App.Repos.Products.CreateProduct(
			t.Context(),
			tx,
			p,
		)

		if err != nil {
			return fmt.Errorf("in create product repo: %w", err)
		}

		return nil

	})
}
