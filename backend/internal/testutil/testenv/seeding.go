package testenv

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
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
