package main

import (
	"log"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/adapters"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/api"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/async"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/config"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

var database *mgo.Database

func main() {
	router := gin.Default()
	config := config.LoadEnvConfig("./.", "app")
	database = util.ConnectToDatabase(config)
	producer := util.ConnectToKafka(config)
	repository := adapters.NewMongoDbStore(database)
	KafkaPublisher := async.NewSyncRequestPublisherServiceImpl(producer, config.SyncRequestTopic)
	service := domain.NewSyncRequestService(repository, KafkaPublisher)
	controller := domain.NewSyncRequestController(service)
	server := api.NewServer(router, &config)
	server.ConfigureController(controller)

	err := server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Failed to srtat HTTP Server")
	}

}
