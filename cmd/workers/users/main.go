package main

import (
	"context"
	"fmt"
	"github.com/s21platform/search-service/internal/databus/user_update"

	kafka_lib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/search-service/internal/clients/user"
	"github.com/s21platform/search-service/internal/config"
	"github.com/s21platform/search-service/internal/databus/user_create"
	"github.com/s21platform/search-service/internal/repository/elastic"
)

func main() {
	cfg := config.MustLoad()
	//logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, "search_kafka", cfg.Platform.Env)
	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "search_kafka", cfg.Platform.Env)
	if err != nil {
		logger_lib.Error(context.Context(context.Background()), fmt.Sprintf("failed to error initializing metrics: %v", err))
		//logger.Error(fmt.Sprintf("failed to error initializing metrics: %v", err))
	}

	userClient := user.MustConnect(cfg)

	elasticClient, err := elastic.MustConnect(cfg)
	if err != nil {
		//logger.Error(fmt.Sprintf("failed to connect to elastic search: %v", err))
		return
	}

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)
	fmt.Println("Start server", cfg.Kafka)

	consumerUserCreateCfg := kafka_lib.DefaultConsumerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.UserCreate, "user.create")
	consumerUserUpdateCfg := kafka_lib.DefaultConsumerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.UserUpdate, "user.update")

	UsersCreateConsumer, err := kafka_lib.NewConsumer(consumerUserCreateCfg, metrics)
	if err != nil {
		//logger.Error(fmt.Sprintf("failed to initialize kafka consumer: %v", err))
		return
	}
	UsersUpdateConsumer, err := kafka_lib.NewConsumer(consumerUserUpdateCfg, metrics)
	if err != nil {
		//logger.Error(fmt.Sprintf("failed to initialize kafka consumer: %v", err))
		return
	}
	newUserCreateHandler := user_create.New(elasticClient, userClient)
	newUserUpdateHandler := user_update.New(elasticClient, userClient)

	UsersCreateConsumer.RegisterHandler(ctx, newUserCreateHandler.Handler)
	UsersUpdateConsumer.RegisterHandler(ctx, newUserUpdateHandler.Handler)

	<-ctx.Done()
}
