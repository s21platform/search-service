package service

import (
	"context"
	"github.com/s21platform/search-service/internal/model"
)

//type userService interface {
//	GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user.GetUserWithOffsetOut, error)
//	CheckFriendship(ctx context.Context, peer string) (bool, error)
//}
//
//type societyService interface {
//	GetSocietyWithOffset(ctx context.Context, limit, offset int64, name string) (*society.GetSocietyWithOffsetOut, error)
//}

type elasticService interface {
	SearchUser(ctx context.Context, query string, filters map[string]interface{}, page, pageSize int) ([]*model.User, int, error)
}
