package users

import (
	"context"

	"github.com/s21platform/search-service/internal/model"

	user_proto "github.com/s21platform/user-proto/user-proto"
)

type Elastic interface {
	Update(ctx context.Context, id string, doc model.UserInfo) error
}

type UserClient interface {
	GetUserInfoByUUID(ctx context.Context, uuid string) (*user_proto.GetUserInfoByUUIDOut, error)
}
