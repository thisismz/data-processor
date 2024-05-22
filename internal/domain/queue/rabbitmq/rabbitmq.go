package rabbitmq

import "github.com/thisismz/data-processor/pkg/rabbitmq"

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
	// TODO : need to fix
	return r.rmq.PublishMessage(data)
}
func (r *RabbitMQRepository) Dequeue() (any, error) {
	// TODO : need to fix
	return r.rmq.ConsumeMessages()
}
