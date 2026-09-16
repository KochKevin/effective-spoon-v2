package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	userssapi "github.com/KochKevin/effective-spoon-v2/internal/users/generated"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Service interface {
	GetUser(ctx context.Context, userId uuid.UUID) (User, error)
}

type Api struct {
	Service Service
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

	userId, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err := a.Service.GetUser(r.Context(), userId)
	if err != nil {
		slog.Error("get user", "err", err)
		render.Status(r, http.StatusInternalServerError)
	}

	render.JSON(w, r, a.ToDto(user))

}
