package auth

import (
	"net/http"

	authapi "github.com/KochKevin/effective-spoon-v2/internal/auth/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
)

type Repo interface {
}

type Api struct {
	Repo Repo
	Txm  infrastructure.TxManager
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/users/generated.ServerInterface

// Login with usercode. Only for testing purposes
// (POST /auth/usercode)
func (a *Api) PostAuthUsercode(w http.ResponseWriter, r *http.Request, params authapi.PostAuthUsercodeParams) {
	panic("not implemented") // TODO: Implement
}
