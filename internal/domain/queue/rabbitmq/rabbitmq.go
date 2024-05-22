package rabbitmq

import (
	"github.com/goccy/go-json"
	"github.com/rabbitmq/amqp091-go"
	"github.com/thisismz/data-processor/pkg/rabbitmq"
)

type RabbitMQRepository struct {
	rmq *rabbitmq.RabbitMQ
}

func New() *RabbitMQRepository {
	rmq := rabbitmq.New()
	return &RabbitMQRepository{
		rmq: rmq,
	}
}

func (r *RabbitMQRepository) Enqueue(data any) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.rmq.PublishMessage(res)
}
func (r *RabbitMQRepository) Dequeue() (<-chan amqp091.Delivery, error) {
	return r.rmq.ConsumeMessages()
}
