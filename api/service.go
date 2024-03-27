package api

import (
	"fmt"

	"github.com/ShapleyIO/iam/api/handlers/identity"
	v1 "github.com/ShapleyIO/iam/api/v1"
	"github.com/ShapleyIO/iam/internal/config"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	*identity.ServiceIdentity
}

var _ v1.ServerInterface = (*Handlers)(nil)

func NewHandlers(cfg *config.Config) (*Handlers, error) {
	handlers := new(Handlers)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",           // no password set
		DB:       cfg.Redis.DB, // use default DB
	})

	handlers.ServiceIdentity = identity.NewServiceIdentity(redisClient)

	return handlers, nil
}
