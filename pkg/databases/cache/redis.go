package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/pkg/env"
)

var (
	REDIS *redis.Client
	ctx   = context.Background()
)

func StartRedis() {
	address := fmt.Sprintf("%s:%s", env.GetEnv("REDIS_HOST", "127.0.0.1"), env.GetEnv("REDIS_PORT", "6379"))
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: env.GetEnv("REDIS_PASSWORD", ""), // no password set
		DB:       0,                                // use default DB
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Err(err).Msg("redis connect failed")
	} else {
		log.Info().Msg("redis connected: " + pong)
		REDIS = client
	}
}
func RedisHealthCheck() bool {
	_, err := REDIS.Ping(ctx).Result()
	if err != nil {
		log.Err(err).Msg("redis health check failed")
	}
	return err == nil
}
func CloseRedis() {
	if REDIS != nil {
		err := REDIS.Close()
		if err != nil {
			log.Err(err).Msg("redis close failed")
		}
	}
	log.Info().Msg("redis closed")
}
