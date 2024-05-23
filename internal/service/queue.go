package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/entity"
)

func SendToQueue() error {
	return queueSrv.queue.Enqueue(entity.User{
		UID:          uuid.New(),
		CreateAt:     time.Now(),
		UserQuota:    "1",
		DataQuota:    "1",
		S3Path:       "s3://bucket-name",
		RateLImit:    0,
		TrafficLImit: 0,
	})
}

func ReceiveFromQueue(done chan bool) error {
	msgs, err := queueSrv.queue.Dequeue()
	if err != nil {
		log.Err(err).Msg("queue dequeue failed")
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			select {
			case <-done:
				// Stop processing messages
				return
			default:
				processMessage(d, done)
			}
		}
	}()

	<-forever
	return nil
}
