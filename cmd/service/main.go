package main

import (
	"fmt"
	"log"
	"net"

	"github.com/s21platform/search-service/internal/clients/society"

	"github.com/s21platform/search-service/internal/clients/friends"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/search-service/internal/clients/user"

	"github.com/s21platform/search-service/internal/infra"

	"github.com/s21platform/search-proto/search"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/rpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	userClient := user.MustConnect(cfg)
	friendsClient := friends.MustConnect(cfg)
	societyClient := society.MustConnect(cfg)
	service := rpc.New(userClient, friendsClient, societyClient)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(infra.Verification),
		grpc.ChainUnaryInterceptor(infra.Logger(logger)),
	)
	search.RegisterSearchServiceServer(server, service)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Println("failed to serve: ", err)
	}
}
