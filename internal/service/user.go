package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/entity"
	"github.com/thisismz/data-processor/pkg/env"
)

var (
	rateLimit           = 0
	trafficLimit        = int64(0)
	ctx                 = context.Background()
	trafficExpiration_i = 0
	rateExpiration_i    = 0
)
var (
	ErrExceededTrafficLimit = errors.New("exceeded traffic limit")
	ErrExceededRateLimit    = errors.New("exceeded rate limit")
	ErrDataDuplicate        = errors.New("data already exists")
)

func GetUser(userQuota string, dataQuota string) (entity.User, error) {
	user, err := storageSrv.store.GetUser(ctx, userQuota)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
func CreateNewUser(userQuota string, dataQuota string) (entity.User, error) {
	nowTime := time.Now()

	trafficExpiration := time.Now().Add(time.Duration(trafficExpiration_i) * time.Minute)
	rateLimitExpiration := time.Now().Add(time.Duration(rateExpiration_i) * time.Minute)

	user := entity.User{
		CreateAt:               nowTime,
		DataQuota:              dataQuota,
		UserQuota:              userQuota,
		UID:                    uuid.New(),
		RateLimit:              rateLimit,
		RateLimitExpiration:    rateLimitExpiration,
		TrafficLimit:           trafficLimit,
		TrafficLimitExpiration: trafficExpiration,
		IsSync:                 true,
	}

	err := storageSrv.store.Add(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func CheckRateLimit(user entity.User) (entity.User, error) {
	nowTime := time.Now()

	rateLimitExpiration := time.Now().Add(time.Duration(rateExpiration_i) * time.Minute)
	if user.RateLimitExpiration.Before(nowTime) {
		if user.RateLimit > 0 {
			user.RateLimit--
		} else {
			return entity.User{}, ErrExceededRateLimit
		}
	} else {
		user.RateLimit = rateLimit
		user.RateLimitExpiration = rateLimitExpiration
	}
	return user, nil
}

func CheckDuplicate(user entity.User) error {
	isDuplicate, err := storageSrv.store.CheckDuplicate(ctx, user.UserQuota, user.DataQuota)
	if err != nil {
		log.Err(err).Msg("check duplicate failed")
		return err
	}
	if isDuplicate {
		return ErrDataDuplicate
	}
	return nil
}

func CheckTrafficLimit(user entity.User, fileSizeBytes int64) (entity.User, error) {
	nowTime := time.Now()

	trafficExpiration := time.Now().Add(time.Duration(trafficExpiration_i) * time.Minute)

	if user.TrafficLimitExpiration.Before(nowTime) {
		if user.TrafficLimit > 0 && user.TrafficLimit >= fileSizeBytes {
			user.TrafficLimit -= fileSizeBytes
		} else {
			return entity.User{}, ErrExceededTrafficLimit
		}
	} else {
		user.TrafficLimit = trafficLimit
		user.TrafficLimitExpiration = trafficExpiration
	}
	return user, nil
}

func UpdateUser(user entity.User) error {
	user.IsSync = true
	user.UserDataQuota = user.UserQuota + ":" + user.DataQuota
	user.CreateAt = time.Now()
	err := storageSrv.store.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	var err error
	rateLimit, err = strconv.Atoi(env.GetEnv("RATE_LIMIT", "1000"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	rateExpiration_i, err = strconv.Atoi(env.GetEnv("RATE_EXPIRATION_MIN", "1"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	trafficLimit, err = strconv.ParseInt(env.GetEnv("TRAFFIC_LIMIT_BYTE", "100000"), 10, 64)
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	trafficExpiration_i, err = strconv.Atoi(env.GetEnv("TRAFFIC_EXPIRATION_MIN", "100"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
}
