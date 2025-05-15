package broker

import (
	"time"

	log "github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker struct {
	connection *amqp.Connection
	config     *Config
}

func NewMessageBroker(config *Config) (*MessageBroker, error) {
	mb := &MessageBroker{
		config: config,
	}

	if err := mb.connect(); err != nil {
		return nil, err
	}

	return mb, nil
}

func (mb *MessageBroker) Close() error {
	if mb.connection != nil {
		return mb.connection.Close()
	}

	return nil
}

func (mb *MessageBroker) DeclareQueue(name string) (*amqp.Queue, error) {
	ch, err := mb.connection.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(name, false, true, true, false, nil)
	if err != nil {
		return nil, err
	}

	return &queue, nil
}

func (mb *MessageBroker) DeclareExchange(name string) error {
	ch, err := mb.connection.Channel()

	if err != nil {
		return err
	}

	return ch.ExchangeDeclare(name, "fanout", true, false, false, false, nil)
}

func (mb *MessageBroker) Publish(exchangeName string, payload []byte) error {
	ch, err := mb.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(exchangeName, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	})
}

func (mb *MessageBroker) Consume(queueName string) (<-chan amqp.Delivery, error) {
	ch, err := mb.connection.Channel()
	if err != nil {
		return nil, err
	}

	return ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (mb *MessageBroker) BindQueue(queueName, exchangeName string) error {
	ch, err := mb.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.QueueBind(
		queueName,
		"",
		exchangeName,
		false,
		nil,
	)
}

func (mb *MessageBroker) connect() error {
	brokerDsn := mb.config.BrokerDsn()

	connection, err := amqp.Dial(brokerDsn)
	if err != nil || connection == nil {
		log.Println("RabbitMQ: waiting to become available")

		for i := 1; i <= 12; i++ {
			connection, err = amqp.Dial(brokerDsn)
			if connection != nil && err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if err != nil || connection == nil {
			return err
		}
	}

	log.Println("RabbitMQ: connection established")

	mb.connection = connection

	return nil
}
