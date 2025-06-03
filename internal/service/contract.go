package service

import (
	"context"

	society "github.com/s21platform/society-proto/society-proto"
	user "github.com/s21platform/user-service/pkg/user"
)

type userService interface {
	GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user.GetUserWithOffsetOut, error)
	CheckFriendship(ctx context.Context, peer string) (bool, error)
}

type societyService interface {
	GetSocietyWithOffset(ctx context.Context, limit, offset int64, name string) (*society.GetSocietyWithOffsetOut, error)
}
