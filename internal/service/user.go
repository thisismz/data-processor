package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/entity"
	"github.com/thisismz/data-processor/pkg/circuit_breaker"
	"github.com/thisismz/data-processor/pkg/env"
)

var (
	rateLimit           = 0
	trafficLimit        = int64(0)
	trafficExpiration   = time.Now()
	rateLimitExpiration = time.Now()
)
var (
	ErrExceededTrafficLimit = errors.New("exceeded traffic limit")
	ErrExceededRateLimit    = errors.New("exceeded rate limit")
	ErrDataDuplicate        = errors.New("data already exists")
)

func UserLimitsCheck(userQuota string, dataQuota string, fileSizeBytes int64) error {
	// check user is exists
	user, err := storageSrv.store.GetUser(context.Background(), userQuota)
	if err != nil {
		return err
	}
	if user.UID == uuid.Nil {
		return createNewUser(userQuota, dataQuota)
	}

	// check rate limit
	if user.RateLimitExpiration.Before(time.Now()) {
		if user.RateLimit > 0 {
			user.RateLimit--
		} else {
			return ErrExceededRateLimit
		}
	} else {
		user.RateLimit = rateLimit
		user.RateLimitExpiration = rateLimitExpiration
	}

	// check data duplicate
	isDuplicate, err := storageSrv.store.CheckDuplicate(context.Background(), userQuota, dataQuota)
	if err != nil {
		log.Err(err).Msg("check duplicate failed")
		return err
	}
	if isDuplicate {
		return ErrDataDuplicate
	}

	// check traffic limit
	if user.TrafficLimitExpiration.Before(time.Now()) {
		if user.TrafficLimit > 0 && user.TrafficLimit >= fileSizeBytes {
			user.TrafficLimit -= fileSizeBytes
		} else {
			return ErrExceededTrafficLimit
		}
	} else {
		user.TrafficLimit = trafficLimit
		user.TrafficLimitExpiration = trafficExpiration
	}

	// update user status
	user.IsSync = true
	user.UserDataQuota = userQuota + ":" + dataQuota
	err = storageSrv.store.Update(context.Background(), user, circuit_breaker.GetCircuitStatus())
	if err != nil {
		return err
	}

	return nil
}
func createNewUser(userQuota string, dataQuota string) error {
	var user entity.User

	user.DataQuota = dataQuota
	user.UserQuota = userQuota
	user.UID = uuid.New()
	user.RateLimit = rateLimit
	user.RateLimitExpiration = rateLimitExpiration
	user.TrafficLimit = trafficLimit
	user.TrafficLimitExpiration = trafficExpiration
	user.IsSync = true

	err := storageSrv.store.Add(context.Background(), user, circuit_breaker.GetCircuitStatus())
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
	trafficLimit, err = strconv.ParseInt(env.GetEnv("TRAFFIC_LIMIT_BYTE", "100000"), 10, 64)
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	trafficExpiration_i, err := strconv.Atoi(env.GetEnv("TRAFFIC_EXPIRATION_MIN", "1"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}
	rateExpiration_i, err := strconv.Atoi(env.GetEnv("RATE_EXPIRATION_MIN", "1"))
	if err != nil {
		log.Err(err).Msg("Error: in converting string to int")
	}

	trafficExpiration = time.Now().Add(time.Duration(trafficExpiration_i) * time.Minute)
	rateLimitExpiration = time.Now().Add(time.Duration(rateExpiration_i) * time.Minute)
}
