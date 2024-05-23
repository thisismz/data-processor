package service

import (
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/internal/entity"
)

func SendToQueue(data entity.Data) error {
	return queueSrv.queue.Enqueue(data)
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
