package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/thisismz/data-processor/internal/domain/queue"
	"github.com/thisismz/data-processor/internal/domain/queue/rabbitmq"
	"github.com/thisismz/data-processor/internal/domain/storage"
	"github.com/thisismz/data-processor/pkg/databases/cache"
	"github.com/thisismz/data-processor/pkg/databases/sql"

	"gorm.io/gorm"
)

type StorageConfiguration func(as *StorageService) error
type QueueConfiquration func(as *QueueService) error

type StorageService struct {
	store storage.StorageRepository
}
type QueueService struct {
	queue queue.QueueRepository
}

var storageSrv *StorageService
var queueSrv *QueueService

func StorageServiceUp() (err error) {
	storageSrv, err = NewStorageService(WithStorageRepository(sql.DataBase, cache.REDIS))
	if err != nil {
		return err
	}
	return nil
}
func QueueServiceUp() (err error) {
	queueSrv, err = NewQueueService(WithQueueRepository())
	if err != nil {
		return err
	}
	return nil
}
func NewStorageService(cfgs ...StorageConfiguration) (*StorageService, error) {
	// Create the orderservice
	os := &StorageService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		if err := cfg(os); err != nil {
			return nil, err
		}
	}
	return os, nil
}
func NewQueueService(cfgs ...QueueConfiquration) (*QueueService, error) {

	os := &QueueService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		if err := cfg(os); err != nil {
			return nil, err
		}
	}
	return os, nil
}

// WithCustomerRepository applies a given customer repository
func WithStorageRepository(dataBaseConnection *gorm.DB, redisClient *redis.Client) StorageConfiguration {
	return func(st *StorageService) error {

		storage := storage.NewStorageRepository(dataBaseConnection, redisClient)
		st.store = storage
		return nil
	}
}

func WithQueueRepository() QueueConfiquration {
	return func(st *QueueService) error {
		// Create the sql repo, if we needed parameters, such as connection strings they could be inputted here
		rmq := rabbitmq.NewRabbitMQRepository()
		st.queue = rmq
		return nil
	}
}
