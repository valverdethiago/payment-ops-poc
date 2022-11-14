package util

import (
	"log"

	"github.com/Pauca-Technologies/payment-ops-poc/sync-manager-ms/config"
	"github.com/Shopify/sarama"
	"gopkg.in/mgo.v2"
)

// LoadEnvConfig load app configuration based on file
func LoadEnvConfig(path string, file string) config.Config {
	config, err := config.LoadConfig(path, file)
	if err != nil {
		log.Fatal("Error loading application config: ", err)
	}
	return config
}

// ConnectToDatabase connect to mongo database
func ConnectToDatabase(config config.Config) *mgo.Database {
	session, err := mgo.Dial(config.DBServer)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	return session.DB(config.DBName)
}

func ConnectToKafka(config config.Config) sarama.SyncProducer {
	configuration := sarama.NewConfig()
	configuration.Producer.Return.Successes = true
	configuration.Producer.RequiredAcks = sarama.WaitForAll
	configuration.Producer.Retry.Max = 5
	brokers := [1]string{config.KafkaBroker}
	con, err := sarama.NewSyncProducer(brokers[:], configuration)
	if err != nil {
		log.Fatal("Error connecting to kafka", err)
	}
	return con
}
