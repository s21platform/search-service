package main

import (
	"context"
	"fmt"

	kafka_lib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/search-service/internal/clients/user"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/databus/users"
	"github.com/s21platform/search-service/internal/repository/elastic"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, "search_kafka", cfg.Platform.Env)
	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "search_kafka", cfg.Platform.Env)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to error initializing metrics: %v", err))
	}

	userClient := user.MustConnect(cfg)

	elasticClient, err := elastic.MustConnect(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to elastic search: %v", err))
		return
	}

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)
	fmt.Println("Start server", cfg.Kafka)

	consumerCfg := kafka_lib.DefaultConsumerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.UserUpdate, "user.update")
	NewUsersConsumer, err := kafka_lib.NewConsumer(consumerCfg, metrics)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initialize kafka consumer: %v", err))
		return
	}
	newUserHandler := users.New(elasticClient, userClient)

	NewUsersConsumer.RegisterHandler(ctx, newUserHandler.Handler)

	<-ctx.Done()
}
