package util

import (
	"database/sql"
	"log"

	"github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/config"
	"github.com/Shopify/sarama"
)

// LoadEnvConfig load app configuration based on file
func LoadEnvConfig(path string, file string) *config.Config {
	config, err := config.LoadConfig(path, file)
	if err != nil {
		log.Fatal("Error loading application config: ", err)
	}
	return &config
}

// ConnectToDatabase connect to mongo database
func ConnectToDatabase(config *config.Config) *sql.DB {
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	return conn
}

func ConnectToKafkaConsumer(config *config.Config) sarama.Consumer {
	configuration := sarama.NewConfig()
	configuration.Consumer.Return.Errors = true
	brokers := [1]string{config.KafkaBroker}
	con, err := sarama.NewConsumer(brokers[:], configuration)
	if err != nil {
		log.Fatal("Error connecting to kafka", err)
	}
	return con
}
