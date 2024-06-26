package service

import (
	"math/rand"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thisismz/data-processor/internal/valueobject"
	"github.com/thisismz/data-processor/pkg/env"
	"github.com/thisismz/data-processor/pkg/hash"
)

func SendDataToQueue(userQuota string, dataQuota string, payload string) error {
	var data valueobject.Data
	// TODO: upload s3 here
	data.UID = uuid.New()
	data.UserQuota = userQuota
	data.DataQuota = dataQuota
	data.Size = int64(len(payload))
	data.Hash = hash.Blake3Hash([]byte(payload))
	data.S3Path = "s3://bucket/path"
	// send data to queue
	err := SendToQueue(data)
	if err != nil {
		return err
	}
	return nil
}
func StartDataProcessor() {
	if env.GetEnv("IS_CONSUMER", "true") == "true" {
		done := make(chan bool)
		go ReceiveFromQueue(done)
		// done
		//done <- true
	} else {
		log.Info().Msg("service act like producer")

	}
}
func processMessage(d amqp091.Delivery, done chan bool) {
	var msg valueobject.Data
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Err(err).Msg("json unmarshal failed")
		d.Nack(false, true) // Requeue message
		return
	}
	log.Info().Msgf("Received data: %v", msg)
	// Simulate processing
	if rand.Intn(2) == 0 { // Simulate random failure
		log.Info().Msgf("Failed to process data: %v", msg)
		d.Nack(false, true) // Requeue message
	} else {
		log.Info().Msgf("Successfully processed data: %v", msg)
		d.Ack(false) // Acknowledge message
	}
	select {
	case <-done:
		// Stop processing messages
		return
	default:
	}
}
