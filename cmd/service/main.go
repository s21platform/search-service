package main

import (
	"fmt"
	"log"
	"net"

	"github.com/s21platform/search-service/internal/clients/user"

	"github.com/s21platform/search-service/internal/infra"

	"github.com/s21platform/search-proto/search"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/rpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	userClient := user.MustConnect(cfg)
	service := rpc.New(userClient)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Verification))
	search.RegisterSearchServiceServer(server, service)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Println("failed to serve: ", err)
	}
}
