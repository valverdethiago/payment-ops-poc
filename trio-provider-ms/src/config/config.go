package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	KafkaBroker            string `mapstructure:"KAFKA_BOOTSTRAP_SERVER"`
	KafkaClientId          string `mapstructure:"TRIO_PROVIDER_KAFKA_CLIENT_ID"`
	DBServer               string `mapstructure:"TRIO_PROVIDER_DB_SERVER"`
	DBName                 string `mapstructure:"TRIO_PROVIDER_DB_NAME"`
	ServerAddress          string `mapstructure:"TRIO_PROVIDER_SERVER_ADDRESS"`
	ReadTimeout            int    `mapstructure:"TRIO_PROVIDER_READ_TIMEOUT"`
	WriteTimeout           int    `mapstructure:"TRIO_PROVIDER_WRITE_TIMEOUT"`
	SyncRequestInputTopic  string `mapstructure:"TRIO_SYNC_REQUEST_INPUT_TOPIC"`
	SyncRequestOutputTopic string `mapstructure:"SYNC_REQUESTS_OUTPUT"`
	BasePath               string `mapstructure:"TRIO_PROVIDER_BASE_PATH"`
	ClientID               string `mapstructure:"TRIO_PROVIDER_CLIENT_ID"`
	ClientSecret           string `mapstructure:"TRIO_PROVIDER_CLIENT_SECRET"`
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
