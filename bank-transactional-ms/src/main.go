package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/api"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/config"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/infra"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/util"
)

var database *sql.DB

func main() {
	config := loadConfig()
	server := configureServer(config)
	startServer(server)
}

func loadConfig() *config.Config {
	return util.LoadEnvConfig(".", "app")
}

func openDatabaseConnection(config *config.Config) *sql.DB {
	return util.ConnectToDatabase(config)
}

func configureService() domain.EventSubscriberService {
	return domain.NewEventSubscriberServiceImpl()
}

func configureConsumer(ctx context.Context, config *config.Config) *infra.Consumer {
	return infra.NewConsumer(ctx, []string{config.KafkaBroker}, config.SyncRequestTopic, "bank-transaction-ms-4")
}

func configureServer(config *config.Config) *api.Server {
	ctx := context.Background()
	server := api.NewServer(config)
	service := configureService()
	consumer := configureConsumer(ctx, config)
	consumer.StartReading(service.OnMessageReceive)
	controller := domain.NewTestController()
	server.ConfigureController(controller)
	return server
}

func startServer(server *api.Server) {
	err := server.Start()
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
