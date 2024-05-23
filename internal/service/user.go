package service

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/domain/storage"
	"github.com/thisismz/data-processor/internal/entity"
	"github.com/thisismz/data-processor/pkg/env"
)

func UserLimitsCheck(userQuota string, dataQuota string) error {
	user, err := storageSrv.store.GetUser(context.Background(), userQuota)
	if err != storage.ErrUserNotFound {
		return err
	}
	if user.UID == uuid.Nil {
		return createNewUser(userQuota, dataQuota)
	}

	return nil
}
func createNewUser(userQuota string, dataQuota string) error {
	var user entity.User

	RateLimit, err := strconv.Atoi(env.GetEnv("RATE_LIMIT", "1000"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	trafficLimit, err := strconv.Atoi(env.GetEnv("TRAFFIC_LIMIT", "100000"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	trafficExpiration_i, err := strconv.Atoi(env.GetEnv("TRAFFIC_EXPIRATION", "1"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	rateExpiration_i, err := strconv.Atoi(env.GetEnv("RATE_EXPIRATION", "1"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}

	trafficExpiration := time.Now().Add(time.Duration(trafficExpiration_i) * time.Minute)
	rateExpiration := time.Now().Add(time.Duration(rateExpiration_i) * time.Minute)

	user.DataQuota = dataQuota
	user.UserQuota = userQuota
	user.UID = uuid.New()
	user.RateLimit = RateLimit
	user.RateLimitExpiration = rateExpiration
	user.TrafficLimit = trafficLimit
	user.TrafficLimitExpiration = trafficExpiration

	err = storageSrv.store.Add(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

func updateUser(user entity.User, data entity.Data) error {
	return storageSrv.store.Update(context.Background(), user)
}
