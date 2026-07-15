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

func (a *AuthService) GetCurrentUserId() uuid.UUID {
	return a.Repo.GetCurrentUserId()
}
