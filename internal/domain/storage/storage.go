package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/domain/storage/memory"
	"github.com/thisismz/data-processor/internal/domain/storage/sql"
	"github.com/thisismz/data-processor/internal/entity"
	"github.com/thisismz/data-processor/pkg/circuit_breaker"
	"gorm.io/gorm"
)

type storageRepository struct {
	memory *memory.RedisRepository
	sql    *sql.SqlRepository
}

func NewStorageRepository(dataBaseConnection *gorm.DB, redisClient *redis.Client) *storageRepository {
	return &storageRepository{
		memory: memory.NewRedisRepository(redisClient),
		sql:    sql.NewSqlRepository(dataBaseConnection),
	}
}
func (st *storageRepository) Add(ctx context.Context, user entity.User) (err error) {
	if circuit_breaker.GetCircuitStatus() {
		if err := st.sql.Create(ctx, user); err != nil {
			log.Err(err).Msg("sql:create user failed")
			return err
		}
		return nil
	}

	if err = st.memory.Add(ctx, user); err != nil {
		user.IsSync = false
		log.Err(err).Msg("redis: unable to set")
	}
	// push to list
	if err = st.memory.Push(ctx, user); err != nil {
		log.Err(err).Msg("redis: unable to push")
	}

	go func() {
		if err := st.sql.Create(ctx, user); err != nil {
			log.Err(err).Msg("sql:create user failed")
		}
	}()
	return err
}

func (st *storageRepository) Update(ctx context.Context, user entity.User) (err error) {

	if circuit_breaker.GetCircuitStatus() {
		if err := st.sql.Create(ctx, user); err != nil {
			log.Err(err).Msg("sql:create user failed")
			return err
		}
		return nil
	}
	if err = st.memory.Update(ctx, user); err != nil {
		user.IsSync = false
		log.Err(err).Msg("redis: unable to update")
	}
	// push to list
	if err = st.memory.Push(ctx, user); err != nil {
		log.Err(err).Msg("redis: unable to push")
	}

	if err = st.memory.SetUserData(ctx, user.UserDataQuota, user.CreateAt); err != nil {
		log.Err(err).Msg("redis: unable to Set")
	}

	go func() {
		if err := st.sql.Create(ctx, user); err != nil {
			log.Err(err).Msg("sql:update user failed")
		}
	}()

	return err
}

func (st *storageRepository) GetData(ctx context.Context, dataQuota string) (entity.User, error) {
	if circuit_breaker.GetCircuitStatus() {
		res, err := st.sql.GetData(ctx, dataQuota)
		if err != nil {
			log.Err(err).Msg("sql:get user failed")
			return res, err
		}
		return res, nil
	}
	data, err := st.memory.GetData(ctx, dataQuota)
	if err != nil {
		if err.Error() == "redis: nil" {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}
	return data, nil
}

func (st *storageRepository) GetUser(ctx context.Context, userQuota string) (entity.User, error) {
	if circuit_breaker.GetCircuitStatus() {
		res, err := st.sql.GetUser(ctx, userQuota)
		if err != nil {
			log.Err(err).Msg("sql:get user failed")
			return res, err
		}
		return res, nil
	}
	user, err := st.memory.GetUser(ctx, userQuota)
	if err != nil {
		if err.Error() == "redis: nil" {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}
	return user, nil
}

func (st *storageRepository) CheckDuplicate(ctx context.Context, userQuota string, dataQuota string) (bool, error) {
	if circuit_breaker.GetCircuitStatus() {
		res, err := st.sql.CheckDuplicate(ctx, userQuota, dataQuota)
		return res, err
	}
	isDuplicate, err := st.memory.CheckDuplicate(ctx, userQuota, dataQuota)
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
	}
	return isDuplicate, nil
}
