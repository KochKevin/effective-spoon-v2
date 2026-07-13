package users

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	userssapi "github.com/KochKevin/effective-spoon-v2/internal/users/generated"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Repo interface {
	GetUser(ctx context.Context, tx *sql.Tx, id uuid.UUID) (User, error)
}

type Api struct {
	Repo Repo
	Txm  infrastructure.TxManager
}

func (a *Api) ToDto(user User) userssapi.User {

	return userssapi.User{
		UserId:  user.Id.String(),
		Name:    user.Name,
		Balance: user.Balance.GetAsEuro(),
	}

}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/users/generated.ServerInterface

// Get the currently logged in user
// (GET /users/current)
func (a *Api) GetUsersCurrent(w http.ResponseWriter, r *http.Request) {

	var dto userssapi.User

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			//slog.Error("Can not get user_id from context")
			return errors.New("Can not get user_id from context")
		}

		cart, err := a.Repo.GetUser(r.Context(), tx, userID)

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting user", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		dto = a.ToDto(cart)

		return nil

	})

	if err != nil {
		slog.Error("Error in /products transaction", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)

}
