package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	KafkaBroker             string `mapstructure:"KAFKA_BOOTSTRAP_SERVER"`
	KafkaClientId           string `mapstructure:"BANK_TRANSACTIONAL_KAFKA_CLIENT_ID"`
	DBSource                string `mapstructure:"BANK_TRANSACTIONAL_DB_SOURCE"`
	DBDriver                string `mapstructure:"BANK_TRANSACTIONAL_DB_DRIVER"`
	ServerAddress           string `mapstructure:"BANK_TRANSACTIONAL_SERVER_ADDRESS"`
	ReadTimeout             int    `mapstructure:"BANK_TRANSACTIONAL_READ_TIMEOUT"`
	WriteTimeout            int    `mapstructure:"BANK_TRANSACTIONAL_WRITE_TIMEOUT"`
	SyncRequestInputTopic   string `mapstructure:"SYNC_REQUEST_INPUT_TOPIC"`
	SyncRequestOutputTopic  string `mapstructure:"SYNC_REQUEST_OUTPUT_TOPIC"`
	BalanceUpdateTopic      string `mapstructure:"BALANCE_UPDATE_TOPIC"`
	TransactionsUpdateTopic string `mapstructure:"TRANSACTIONS_UPDATE_TOPIC"`
}

// LoadConfig loads config from env
func LoadConfig(path string, env string) (Config, error) {
	var config Config
	configPath := fmt.Sprintf("../%s/env", path)
	viper.AddConfigPath(configPath)
	viper.SetConfigName(env)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	return config, err
}

// LoadEnvConfig load app configuration based on file
func LoadEnvConfig(path string, file string) Config {
	config, err := LoadConfig(path, file)
	if err != nil {
		log.Fatal("Error loading application config: ", err)
	}
	return config
}
