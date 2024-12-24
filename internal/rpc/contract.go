package rpc

import (
	"context"

	user_proto "github.com/s21platform/user-proto/user-proto"
)

type userService interface {
	GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user_proto.GetUserWithOffsetOut, error)
}

type friendsService interface {
	IsFriendsExist(ctx context.Context, uuid string) (bool, error)
}
