package main

import (
	"fmt"
	"github.com/s21platform/search-service/internal/infra"
	"log"
	"net"

	"github.com/s21platform/search-proto/search"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/rpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	service := rpc.New()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.UnaryInterceptor))
	search.RegisterSearchServiceServer(server, service)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Println("failed to serve: ", err)
	}
}
