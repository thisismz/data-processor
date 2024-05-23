package service

import (
	"math/rand"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thisismz/data-processor/internal/entity"
)

func processMessage(d amqp091.Delivery, done chan bool) {
	var msg entity.User
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
