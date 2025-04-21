package user

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	user_proto "github.com/s21platform/user-proto/user-proto"

	"github.com/s21platform/search-service/internal/config"
)

type Client struct {
	client user_proto.UserServiceClient
}

func (c *Client) GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user_proto.GetUserWithOffsetOut, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	users, err := c.client.GetUserWithOffset(ctx, &user_proto.GetUserWithOffsetIn{Limit: limit, Offset: offset, Nickname: nickName})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return users, nil
}

func MustConnect(cfg *config.Config) *Client {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.User.Host, cfg.User.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	client := user_proto.NewUserServiceClient(conn)
	return &Client{client: client}
}

func (c *Client) GetUserInfoByUUID(ctx context.Context, uuid string) (*user_proto.GetUserInfoByUUIDOut, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	userInfo, err := c.client.GetUserInfoByUUID(ctx, &user_proto.GetUserInfoByUUIDIn{Uuid: uuid})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return userInfo, nil
}

func (c *Client) GetUsersInfoWithOffset(ctx context.Context, nickname string, limit, total int64) (*user_proto.GetUserWithOffsetOutAll, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	usersInfo, err := c.client.GetUsersInfoWithOffset(ctx, &user_proto.GetUserWithOffsetIn{Limit: limit, Offset: total, Nickname: nickname})
	if err != nil {
		return nil, fmt.Errorf("failet to get all users: %w", err)
	}
	return usersInfo, nil
}
