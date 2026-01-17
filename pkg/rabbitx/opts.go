package rabbitx

import amqp "github.com/rabbitmq/amqp091-go"

type Option func(m *amqp.Publishing)

func WithPersistent() Option {
	return func(m *amqp.Publishing) {
		m.DeliveryMode = amqp.Persistent
	}
}

func WithContentType(contentType string) Option {
	return func(m *amqp.Publishing) {
		m.ContentType = contentType
	}
}

func WithMessageIdGenerator(messageIdGenerator func() string) Option {
	return func(m *amqp.Publishing) {
		m.MessageId = messageIdGenerator()
	}
}
