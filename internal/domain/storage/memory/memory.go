package memory

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/entity"
	"github.com/thisismz/data-processor/pkg/env"
)

var listName = "date"

type RedisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(redis *redis.Client) *RedisRepository {
	return &RedisRepository{
		redis: redis,
	}
}

func (r *RedisRepository) Add(ctx context.Context, user entity.User) error {
	if err := r.redis.Set(ctx, user.UserQuota, user, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) SetUserData(ctx context.Context, key string, val any) error {

	timeLimit, err := strconv.Atoi(env.GetEnv("DUPLICATE_TIME_LIMIT", "10"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}

	if err := r.redis.Set(ctx, key, val, time.Duration(timeLimit)*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) GetUser(ctx context.Context, userQuota string) (entity.User, error) {
	var user entity.User
	if err := r.redis.Get(ctx, userQuota).Scan(&user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *RedisRepository) GetData(ctx context.Context, dataQuota string) (entity.User, error) {
	var user entity.User
	if err := r.redis.Get(ctx, dataQuota).Scan(&user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *RedisRepository) FlushRedis() error {
	if err := r.redis.FlushAll(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}

// push to redis
func (r *RedisRepository) Push(ctx context.Context, user entity.User) error {
	if err := r.redis.RPush(ctx, listName, user); err != nil {
		return err.Err()
	}
	return nil
}

// pull from redis last one
func (r *RedisRepository) Pull(ctx context.Context) (entity.User, error) {
	var user entity.User
	if err := r.redis.LPop(ctx, listName).Scan(&user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// remove all from redis list
func (r *RedisRepository) RemoveAllFromList(ctx context.Context) error {
	if err := r.redis.Del(ctx, listName).Err(); err != nil {
		return err
	}
	return nil
}
func (r *RedisRepository) Update(ctx context.Context, user entity.User) error {
	if err := r.redis.Set(ctx, user.UserQuota, user, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) CheckDuplicate(ctx context.Context, userQuota string, dataQuota string) (bool, error) {
	key := userQuota + ":" + dataQuota
	res, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if res == "" {
		return false, nil
	}
	return true, nil
}
