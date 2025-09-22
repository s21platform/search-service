package main

import (
	"context"
	"fmt"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/search-service/internal/repository/elastic"
	"log"
	"net"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/search-service/internal/infra"

	"github.com/s21platform/search-proto/search"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, "search_kafka", cfg.Platform.Env)
	ctx := logger_lib.NewContext(context.Background(), logger)
	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "search", cfg.Platform.Env)
	_ = metrics
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to create metrics object: %v", err))
		log.Fatalf("failed to create metrics object: %v", err)
	}

	elasticClient, err := elastic.MustConnect(cfg)
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to connect to elastic search: %v", err))
		return
	}

	serviceElastic := service.New(elasticClient)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(infra.Verification),
		grpc.ChainUnaryInterceptor(infra.Logger(logger)),
	)
	search.RegisterSearchServiceServer(server, serviceElastic)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Println("failed to serve: ", err)
	}
}
