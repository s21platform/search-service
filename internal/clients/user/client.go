package user

import (
	"context"
	"fmt"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/s21platform/search-service/internal/config"
	user "github.com/s21platform/user-service/pkg/user"
)

type Client struct {
	client user.UserServiceClient
}

func (c *Client) GetUserWithOffset(ctx context.Context, limit, offset int64, nickName string) (*user.GetUserWithOffsetOut, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	users, err := c.client.GetUserWithOffset(ctx, &user.GetUserWithOffsetIn{Limit: limit, Offset: offset, Nickname: nickName})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return users, nil
}

func (c *Client) CheckFriendship(ctx context.Context, peer string) (bool, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	isFriend, err := c.client.CheckFriendship(ctx, &user.CheckFriendshipIn{
		Uuid: peer,
	})
	if err != nil {
		return false, err
	}
	return isFriend.Succses, nil
}

func MustConnect(cfg *config.Config) *Client {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.User.Host, cfg.User.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	client := user.NewUserServiceClient(conn)
	return &Client{client: client}
}
