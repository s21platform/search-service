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
	GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user_proto.GetUserWithOffsetOut, error)
	GetUsersInfoWithOffset(ctx context.Context, nickname string, limit, total int64) (*user_proto.GetUserWithOffsetOutAll, error)
}
