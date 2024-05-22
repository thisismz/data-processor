package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/thisismz/data-processor/pkg/env"
)

type RabbitMQ struct {
	conn         *amqp091.Connection
	channel      *amqp091.Channel
	queueName    string
	exchangeName string
	routingKey   string
}

func New() *RabbitMQ {
	exchangeName := env.GetEnv("RABBITMQ_EXCHANGE_NAME", "direct_exchange")
	queueName := env.GetEnv("RABBITMQ_QUEUE_NAME", "worker_queue")
	routingKey := env.GetEnv("RABBITMQ_ROUTING_KEY", "worker_key")

	conn, err := connectToRabbitMQ()
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := createChannel(conn)
	failOnError(err, "Failed to open a channel")

	err = declareExchange(channel, exchangeName)
	failOnError(err, "Failed to declare an exchange")

	queue, err := declareQueue(channel, queueName)
	failOnError(err, "Failed to declare a queue")

	err = bindQueue(channel, queue.Name, exchangeName, routingKey)
	failOnError(err, "Failed to bind a queue")
	return &RabbitMQ{
		conn:         conn,
		channel:      channel,
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,
	}
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Err(err).Msg(msg)
	}
}

func connectToRabbitMQ() (*amqp091.Connection, error) {
	conn, err := amqp091.Dial(env.GetEnv("RABBITMQ_URL", "amqp:///"))
	return conn, err
}

func createChannel(conn *amqp091.Connection) (*amqp091.Channel, error) {
	channel, err := conn.Channel()
	return channel, err
}
func declareExchange(channel *amqp091.Channel, exchangeName string) error {
	err := channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err
}
func declareQueue(channel *amqp091.Channel, queueName string) (amqp091.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return queue, err
}
func bindQueue(channel *amqp091.Channel, queueName, exchangeName, routingKey string) error {
	err := channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	return err
}
func (r *RabbitMQ) PublishMessage(message []byte) error {
	err := r.channel.Publish(
		r.exchangeName, // exchange
		r.routingKey,   // routing key
		false,          // mandatory
		false,          // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	return err
}
func (r *RabbitMQ) ConsumeMessages() (<-chan amqp091.Delivery, error) {
	msgs, err := r.channel.Consume(
		r.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	return msgs, err
}

func (r *RabbitMQ) Close() {
	if err := r.channel.Close(); err != nil {
		log.Err(err).Msg("rabbitmq close failed")
	}

	if err := r.conn.Close(); err != nil {
		log.Err(err).Msg("rabbitmq close failed")
	}

	log.Info().Msg("rabbitmq connection closed")
}
