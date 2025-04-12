package users

import (
	"context"

	user_proto "github.com/s21platform/user-proto/user-proto"
)

type Elastic interface{}

type UserClient interface {
	GetUserInfoByUUID(ctx context.Context, uuid string) (*user_proto.GetUserInfoByUUIDOut, error)
}
