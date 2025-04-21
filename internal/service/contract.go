package service

import (
	"context"

	"github.com/s21platform/search-service/internal/model"

	society "github.com/s21platform/society-proto/society-proto"

	user_proto "github.com/s21platform/user-proto/user-proto"
)

type userService interface {
	GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user_proto.GetUserWithOffsetOut, error)
	GetUsersInfoWithOffset(ctx context.Context, nickname string, limit, total int64) (*user_proto.GetUserWithOffsetOutAll, error)
}

type friendsService interface {
	IsFriendsExist(ctx context.Context, uuid string) (bool, error)
}

type societyService interface {
	GetSocietyWithOffset(ctx context.Context, limit, offset int64, name string) (*society.GetSocietyWithOffsetOut, error)
}

type Elastic interface {
	Update(ctx context.Context, id string, doc model.UserInfo) error
	ExistOrCreateIndex(ctx context.Context, users []string) error
	BulkIndexUsers(ctx context.Context, users []model.UserInfo) error
}
