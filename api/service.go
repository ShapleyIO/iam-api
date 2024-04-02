package api

import (
	"fmt"

	"github.com/ShapleyIO/iam/api/handlers/authn"
	"github.com/ShapleyIO/iam/api/handlers/health"
	"github.com/ShapleyIO/iam/api/handlers/identity"
	v1 "github.com/ShapleyIO/iam/api/v1"
	"github.com/ShapleyIO/iam/internal/config"
	"github.com/ShapleyIO/iam/internal/passwordhasher"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	*identity.ServiceIdentity
	*authn.ServiceAuthN
	*health.ServiceHealth
}

var _ v1.ServerInterface = (*Handlers)(nil)

func NewHandlers(cfg *config.Config) (*Handlers, error) {
	handlers := new(Handlers)

	// Create a Redis Client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",           // no password set
		DB:       cfg.Redis.DB, // use default DB
	})

	// Create Password Hasher
	ph := passwordhasher.NewArgonHasher(cfg)

	handlers.ServiceIdentity = identity.NewServiceIdentity(redisClient, ph)
	handlers.ServiceAuthN = authn.NewServiceAuthN(redisClient, ph)
	handlers.ServiceHealth = health.NewServiceHealth(redisClient)

	return handlers, nil
}
