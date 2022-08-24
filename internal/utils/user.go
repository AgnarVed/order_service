package utils

import (
	"context"
	"errors"
	"tests2/internal/models"
)

type UtilsUser interface {
	GetUserFromContext(ctx context.Context) (*models.UserFromRequest, error)
	PutUserToContext(ctx context.Context, user *models.UserFromRequest) (context.Context, error)
}

var ErrFile = errors.New("Некорректный тип файла.")

type utilsUser struct{}

type keyCtx int

const (
	keyCtxUser keyCtx = iota
)

func (u *utilsUser) GetUserFromContext(ctx context.Context) (*models.UserFromRequest, error) {
	userCtx := ctx.Value(keyCtxUser)
	if userCtx == nil {
		return nil, errors.New("can't get user from context")
	}

	user := userCtx.(*models.UserFromRequest)
	return user, nil
}
func (u *utilsUser) PutUserToContext(ctx context.Context, user *models.UserFromRequest) (context.Context, error) {
	ctxBack := context.WithValue(ctx, keyCtxUser, user)

	return ctxBack, nil
}

func NewUtilsUser() UtilsUser {
	return &utilsUser{}
}
