package circuit_breaker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/domain/storage/memory"
	sqlRep "github.com/thisismz/data-processor/internal/domain/storage/sql"
	"github.com/thisismz/data-processor/pkg/databases/cache"
	"github.com/thisismz/data-processor/pkg/databases/sql"
)

func sync() {
	ctx := context.Background()
	if IsSync() {
		cache.RedisSyncStatus = true
		return
	}
	memory := memory.NewRedisRepository(cache.REDIS)
	lastDate, err := memory.Pull(ctx)
	if err != nil {
		log.Err(err).Msg("memory:pull failed")
		return
	}
	// get users
	sql := sqlRep.NewSqlRepository(sql.DataBase)
	listOfUsers, err := sql.GetSync(ctx, lastDate.CreateAt)
	if err != nil {
		log.Err(err).Msg("sql:get sync failed")
	}
	log.Info().Msgf("list of users: %v", listOfUsers)
	if len(listOfUsers) <= 0 {
		return
	}

	// sync users
	for _, user := range listOfUsers {
		// sub create date from now date
		if user.TrafficLimitExpiration.Compare(time.Now()) >= 0 {
			log.Info().Msg("user is expired")
			continue
		}
		user.IsSync = true

		if err := memory.Update(ctx, user); err != nil {
			log.Err(err).Msg("memory:add failed")
			break
		}

		if err := memory.SetUserData(ctx, user.UserDataQuota, user.CreateAt); err != nil {
			log.Err(err).Msg("memory:add failed")
		}

		err := sql.Update(ctx, user)
		if err != nil {
			log.Err(err).Msg("sql:update failed")
			break
		}
	}
}
func IsSync() bool {
	sql := sqlRep.NewSqlRepository(sql.DataBase)
	listOfUsers, err := sql.GetUnSync()
	if err != nil {
		log.Err(err).Msg("sql:get sync failed")
	}
	if len(listOfUsers) <= 0 {
		return true
	}
	return false
}
