package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Database       *DatabaseConfig
	Redis          *RedisConfig
	PasswordHasher *PasswordHasherConfig
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

type PasswordHasherConfig struct {
	SaltLength  int
	KeyLength   uint32
	Time        uint32
	Memory      uint32
	Threads     uint8
	HashVariant string
}

func NewConfig() (config *Config, err error) {
	viper.AutomaticEnv()

	// Map Viper Variables to ENV Aliases
	viper.RegisterAlias("Database.Host", "IAM_DATABASE_HOST")
	viper.RegisterAlias("Database.Port", "IAM_DATABASE_PORT")
	viper.RegisterAlias("Database.User", "IAM_DATABASE_USER")

	viper.RegisterAlias("Redis.Host", "IAM_REDIS_HOST")
	viper.RegisterAlias("Redis.Port", "IAM_REDIS_PORT")
	viper.RegisterAlias("Redis.DB", "IAM_REDIS_DB")

	viper.RegisterAlias("PasswordHasher.SaltLength", "IAM_PASSWORD_HASHER_SALT_LENGTH")
	viper.RegisterAlias("PasswordHasher.KeyLength", "IAM_PASSWORD_HASHER_KEY_LENGTH")
	viper.RegisterAlias("PasswordHasher.Time", "IAM_PASSWORD_HASHER_TIME")
	viper.RegisterAlias("PasswordHasher.Memory", "IAM_PASSWORD_HASHER_MEMORY")
	viper.RegisterAlias("PasswordHasher.Threads", "IAM_PASSWORD_HASHER_THREADS")
	viper.RegisterAlias("PasswordHasher.HashVariant", "IAM_PASSWORD_HASHER_HASH_VARIANT")

	// Set Defaults to Local Env
	viper.SetDefault("Database.Host", "localhost")
	viper.SetDefault("Database.Port", "5432")
	viper.SetDefault("Database.User", "postgres")

	viper.SetDefault("Redis.Host", "localhost")
	viper.SetDefault("Redis.Port", "6379")
	viper.SetDefault("Redis.DB", 0)

	viper.SetDefault("PasswordHasher.SaltLength", 24)
	viper.SetDefault("PasswordHasher.KeyLength", 32)
	viper.SetDefault("PasswordHasher.Time", 4)
	viper.SetDefault("PasswordHasher.Memory", 128*1024)
	viper.SetDefault("PasswordHasher.Threads", 8)
	viper.SetDefault("PasswordHasher.HashVariant", "argon2id")

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	log.Info().Interface("config", config).Msg("loaded config")

	return config, nil
}
