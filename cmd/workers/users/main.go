package main

import (
	"context"
	"fmt"

	"github.com/s21platform/search-service/internal/clients/user"

	kafka_lib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/databus/users"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, "search_kafka", cfg.Platform.Env)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "search_kafka", cfg.Platform.Env)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to error initializing metrics: %v", err))
	}

	userClient := user.MustConnect(cfg)

	// worker elastic
	var elastic config.Elastic
	// end init elastic

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)
	fmt.Println("Start server", cfg.Kafka.Server, cfg.Kafka.UserUpdate)

	NewUsersConsumer, err := kafka_lib.NewConsumer(cfg.Kafka.Server, cfg.Kafka.UserUpdate, metrics)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initialize kafka consumer: %v", err))
		return
	}

	newUserHandler := users.New(elastic, userClient)

	NewUsersConsumer.RegisterHandler(ctx, newUserHandler.Handler)

	<-ctx.Done()
}
