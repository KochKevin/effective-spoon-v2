package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	authapi "github.com/KochKevin/effective-spoon-v2/internal/auth/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserIdByCode(ctx context.Context, tx *sql.Tx, usercode users.Usercode) (uuid.UUID, error)
}

type AuthService interface {
	Login(userId uuid.UUID) error
	Logout()
}

type Api struct {
	UserRepo    UserRepo
	AuthService AuthService
	Txm         infrastructure.TxManager
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/auth/generated.ServerInterface

//LOGIN
//curl -X POST "https://effective-waddle-4j7xj7vp44pxh7xrp-8080.app.github.dev/api/auth/usercode?usercode=3"

//LOGOUT
//curl -X POST "https://effective-waddle-4j7xj7vp44pxh7xrp-8080.app.github.dev/api/auth/logout"


// Login with usercode. Only for testing purposes
// (POST /auth/usercode)
func (a *Api) PostAuthUsercode(w http.ResponseWriter, r *http.Request, params authapi.PostAuthUsercodeParams) {

	slog.Debug("Try to login with ", params.Usercode, " usercode")

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		userId, err := a.UserRepo.GetUserIdByCode(r.Context(), tx, users.Usercode(params.Usercode))

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting user id by user code", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		slog.Debug("Found userId with the ", "usercode", params.Usercode, " userId: ", userId)

		err = a.AuthService.Login(userId)
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error trying to login user in Auth Service", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}


		return nil

	})

	if err != nil {
		slog.Error("Error in login with usercode transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
}

// Logout current user
// (POST /auth/logout)
func (a *Api) PostAuthLogout(w http.ResponseWriter, r *http.Request) {

	a.AuthService.Logout()

	render.Status(r, http.StatusOK)

}
