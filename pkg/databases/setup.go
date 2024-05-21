package databases

import (
	"github.com/thisismz/data-processor/pkg/databases/cache"
	"github.com/thisismz/data-processor/pkg/databases/sql"
)

func StartDatabase() {
	cache.StartRedis()
	sql.StartMysql()
}

func CloseDatabase() {
	cache.CloseRedis()
	sql.CloseMysql()
}
