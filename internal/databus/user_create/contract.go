package user_create

import (
	"context"
	"github.com/s21platform/search-service/internal/model"
)

type Elastic interface {
	SetUser(ctx context.Context, user *model.User) error
}

type UserClient interface {
}
