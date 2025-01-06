package society

import (
	"context"
	"fmt"
	"log"

	"github.com/s21platform/search-service/internal/config"
	society "github.com/s21platform/society-proto/society-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client society.SocietyServiceClient
}

func MustConnect(cfg *config.Config) *Client {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Society.Host, cfg.Society.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to society: %v", err)
	}
	client := society.NewSocietyServiceClient(conn)
	return &Client{client: client}
}

func (c *Client) GetSocietyWithOffset(ctx context.Context, limit, offset int64, name string) (*society.GetSocietyWithOffsetOut, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
	societies, err := c.client.GetSocietyWithOffset(ctx, &society.GetSocietyWithOffsetIn{
		Offset: offset,
		Limit:  limit,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}
	return societies, nil
}
