package queue

import "github.com/rabbitmq/amqp091-go"

type QueueRepository interface {
	Enqueue(any) error
	Dequeue() (<-chan amqp091.Delivery, error)
}
