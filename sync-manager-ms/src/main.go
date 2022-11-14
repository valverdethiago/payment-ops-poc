package main

import (
	"context"
	"log"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/adapters"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/api"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/config"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/infra"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/util"
	"gopkg.in/mgo.v2"
)

var database *mgo.Database

func main() {
	config := loadConfig()
	server := configureServer(config)
	startServer(config, server)
}

func loadConfig() *config.Config {
	return util.LoadEnvConfig(".", "app")
}

func openDatabaseConnection(config *config.Config) *mgo.Database {
	return util.ConnectToDatabase(config)
}

func configureKafkaProducer(ctx context.Context, config *config.Config) *infra.Producer {
	return infra.NewProducer(ctx, []string{config.KafkaBroker}, config.SyncRequestTopic)
}

func configureService(config *config.Config, producer *infra.Producer) domain.SyncRequestService {
	database := openDatabaseConnection(config)
	repository := adapters.NewMongoDbStore(database)
	service := domain.NewSyncRequestService(repository, producer)
	return service
}

func configureServer(config *config.Config) *api.Server {
	ctx := context.Background()
	server := api.NewServer(config)
	producer := configureKafkaProducer(ctx, config)
	service := configureService(config, producer)
	controller := domain.NewSyncRequestController(service)
	server.ConfigureController(controller)
	return server
}

func startServer(config *config.Config, server *api.Server) {
	err := server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
