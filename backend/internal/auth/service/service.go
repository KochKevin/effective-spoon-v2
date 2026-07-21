package authservice

import (
	"errors"

	"github.com/google/uuid"
)

type Repo interface {
	SetCurrentUserId(userId uuid.UUID)
	GetCurrentUserId() (userId uuid.UUID)
	TryLogin(userId uuid.UUID) bool
}

type AuthService struct {
	Repo Repo
}

func (a *AuthService) Login(userId uuid.UUID) error {

	if a.Repo.TryLogin(userId) {
		return nil
	}

	return errors.New("Error: Login is not possible now. An other user is already logged in.")
}

func (a *AuthService) Logout() {

	a.Repo.SetCurrentUserId(uuid.Nil)

}

var NoCurrentUser = errors.New("error: no current user is setted to be getted. Set an current user first")

func (a *AuthService) GetCurrentUserId() (uuid.UUID, error) {

	userId := a.Repo.GetCurrentUserId()

	if userId == uuid.Nil {
		return uuid.Nil, NoCurrentUser
	}

	return userId, nil
}
