package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RabbitDSN               string `yaml:"rabbit_dsn" envconfig:"RABBIT_DSN" default:"" required:"true"`
	DatabaseDSN             string `yaml:"database_dsn" envconfig:"DATABASE_DSN" default:"" required:"true"`
	DBMaxOpenConns          int    `yaml:"db_max_open_conns" envconfig:"DB_MAX_OPEN_CONNS" default:"25" required:"true"`
	DBMaxIdleConns          int    `yaml:"db_max_idle_conns" envconfig:"DB_MAX_IDLE_CONNS" default:"25" required:"true"`
	DBConnMaxLifetime       int    `yaml:"db_conn_max_lifetime" envconfig:"DB_CONN_MAX_LIFETIME" default:"5" required:"true"`
	RabbitReconnectInterval int    `yaml:"rabbit_reconnect_interval" envconfig:"RABBIT_RECONNECT_INTERVAL" default:"5" required:"true"`
	RabbitGoroutinesCnt     int    `yaml:"rabbit_goroutines_cnt" envconfig:"RABBIT_GOROUTINES_CNT" default:"10" required:"true"`
	TableName               string `yaml:"table_name" envconfig:"TABLE_NAME" default:"messages" required:"true"`
	Queue                   string `yaml:"queue" envconfig:"QUEUE" default:"" required:"true"`
	ConsumerName            string `yaml:"consumer_name" envconfig:"CONSUMER_NAME" default:"go-consumer" required:"true"`
	LogFormat               string `yaml:"log_format" envconfig:"LOG_FORMAT" default:"text" required:"true"`
}

// GetConfig gets path to yaml file. If path is an empty string,
// configuration will be obtained from environment variables
func GetConfig(path string) (Config, error) {
	if path == "" {
		return getConfigFromEnv()
	}

	return getConfigFromYaml(path)
}

func getConfigFromYaml(path string) (Config, error) {
	config := Config{
		RabbitDSN:               "",
		DatabaseDSN:             "",
		DBMaxOpenConns:          0,
		DBMaxIdleConns:          0,
		DBConnMaxLifetime:       0,
		RabbitReconnectInterval: 0,
		RabbitGoroutinesCnt:     0,
		Queue:                   "",
		ConsumerName:            "",
		LogFormat:               "",
		TableName:               "",
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	err = yaml.NewDecoder(file).Decode(&config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func getConfigFromEnv() (Config, error) {
	config := Config{
		RabbitDSN:               "",
		DatabaseDSN:             "",
		DBMaxOpenConns:          0,
		DBMaxIdleConns:          0,
		DBConnMaxLifetime:       0,
		RabbitReconnectInterval: 0,
		RabbitGoroutinesCnt:     0,
		Queue:                   "",
		ConsumerName:            "",
		LogFormat:               "",
		TableName:               "",
	}

	err := envconfig.Process("", &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
