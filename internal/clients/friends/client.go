package friends

import (
	"context"
	"fmt"
	"log"

	friends_proto "github.com/s21platform/friends-proto/friends-proto"
	"github.com/s21platform/search-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client friends_proto.FriendsServiceClient
}

func MustConnect(cfg *config.Config) *Client {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Friends.Host, cfg.Friends.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	client := friends_proto.NewFriendsServiceClient(conn)
	return &Client{client: client}
}

func (c *Client) IsFriendsExist(ctx context.Context, uuid string) (bool, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	isFriend, err := c.client.IsFriendExist(ctx, &friends_proto.IsFriendExistIn{
		Peer: uuid,
	})
	if err != nil {
		return false, err
	}
	return isFriend.Success, nil
}
