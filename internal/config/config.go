package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database *DatabaseConfig
	Redis    *RedisConfig
}

type DatabaseConfig struct {
	Host string
	Port int
	User string
}

type RedisConfig struct {
	Host string
	Port int
	DB   int
}

func NewConfig() (config *Config, err error) {
	viper.AutomaticEnv()

	// Set Defaults to Local Env
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
