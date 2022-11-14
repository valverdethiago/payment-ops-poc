package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	KafkaBroker   string `mapstructure:"KAFKA_BOOTSTRAP_SERVER"`
	kafkaClientId string `mapstructure:"BANK_TRANSACTIONAL_KAFKA_CLIENT_ID"`
	DBSource      string `mapstructure:"BANK_TRANSACTIONAL_DB_SOURCE"`
	DBDriver      string `mapstructure:"BANK_TRANSACTIONAL_DB_DRIVER"`
	ServerAddress string `mapstructure:"BANK_TRANSACTIONAL_SERVER_ADDRESS"`
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
