package authservice

import (
	"context"
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

	return errors.New("login is not possible now. An other user is already logged in")
}

func (a *AuthService) Logout() {

	a.Repo.SetCurrentUserId(uuid.Nil)

}

var ErrNoCurrentUser = errors.New("error: no current user is setted to be getted. Set an current user first")

func (a *AuthService) GetCurrentUserId() (uuid.UUID, error) {

	userId := a.Repo.GetCurrentUserId()

	if userId == uuid.Nil {
		return uuid.Nil, ErrNoCurrentUser
	}

	return userId, nil
}

type ctxKey string

const UserIdCtxKey ctxKey = "user_id"

func WithUser(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIdCtxKey, userId)
}

func UserIdFrom(ctx context.Context) (userId uuid.UUID, ok bool) {
	userId, ok = ctx.Value(UserIdCtxKey).(uuid.UUID)
	return userId, ok
}
