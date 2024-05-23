package memory

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thisismz/data-processor/internal/entity"
)

type RedisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(redis *redis.Client) *RedisRepository {
	return &RedisRepository{
		redis: redis,
	}
}

func (r *RedisRepository) Add(ctx context.Context, user entity.User, expiration time.Duration) error {
	if err := r.redis.Set(ctx, user.UserQuota, user, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) AddUserDataKey(ctx context.Context, user entity.User) error {
	key := user.UserQuota + ":" + user.DataQuota
	if err := r.redis.Set(ctx, key, user, 0).Err(); err != nil {
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

func (r *RedisRepository) Update(ctx context.Context, user entity.User) error {
	if err := r.redis.Set(ctx, user.UserQuota, user, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) CheckDuplicate(ctx context.Context, userQuota string, dataQuota string) (bool, error) {
	key := userQuota + ":" + dataQuota
	return r.redis.Get(ctx, key).Bool()
}
