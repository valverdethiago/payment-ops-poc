package main

import (
	"log"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/adapters"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/api"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/async"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/config"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/util"
	"github.com/Shopify/sarama"
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

func connectToKafka(config *config.Config) sarama.SyncProducer {
	return util.ConnectToKafka(config)
}

func configureService(config *config.Config, producer sarama.SyncProducer) domain.SyncRequestService {
	database := openDatabaseConnection(config)
	repository := adapters.NewMongoDbStore(database)
	KafkaPublisher := async.NewSyncRequestPublisherServiceImpl(producer, config.SyncRequestTopic)
	service := domain.NewSyncRequestService(repository, KafkaPublisher)
	return service
}

func configureServer(config *config.Config) *api.Server {
	server := api.NewServer(config)
	producer := connectToKafka(config)
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
