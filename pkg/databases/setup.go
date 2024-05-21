package databases

import "github.com/thisismz/data-processor/pkg/databases/cache"

func StartDatabase() {
	cache.StartRedis()
}

func CloseDatabase() {
	cache.CloseRedis()
}
