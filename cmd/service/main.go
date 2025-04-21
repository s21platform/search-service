package main

import (
	"fmt"
	"log"
	"net"

	"github.com/s21platform/search-service/internal/repository/elsearch"

	"github.com/s21platform/search-service/internal/clients/society"

	"github.com/s21platform/search-service/internal/clients/friends"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/search-service/internal/clients/user"

	"github.com/s21platform/search-service/internal/infra"

	"github.com/s21platform/search-proto/search"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	fmt.Println("start \n\n cfg user:", cfg.User.Host, cfg.User.Port)
	userClient := user.MustConnect(cfg)
	friendsClient := friends.MustConnect(cfg)
	societyClient := society.MustConnect(cfg)
	elastic, err := elsearch.New(cfg.Elastic)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create elastic client: %v", err))
		return
	}
	service := service.New(userClient, elastic, friendsClient, societyClient)

	//load all users
	if err := service.LoadAllUsers(); err != nil {
		fmt.Println("error load", err)
	} else {
		fmt.Println("loaded all users")
	}

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
