package user

import (
	"context"
	"fmt"
	"log"

	"github.com/s21platform/search-service/internal/config"
	user_proto "github.com/s21platform/user-proto/user-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handle struct {
	client user_proto.UserServiceClient
}

func (h *Handle) GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user_proto.GetUserWithOffsetOut, error) {
	users, err := h.client.GetUserWithOffset(ctx, &user_proto.GetUserWithOffsetIn{Limit: limit, Offset: offset, Nickname: nickName})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return users, nil
}

func MustConnect(cfg *config.Config) *Handle {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.User.Host, cfg.User.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	client := user_proto.NewUserServiceClient(conn)
	return &Handle{client: client}
}
