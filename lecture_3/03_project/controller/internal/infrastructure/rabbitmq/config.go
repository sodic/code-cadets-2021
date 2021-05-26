package rabbitmq

import "github.com/streadway/amqp"

// ConsumerConfig determines how a RabbitMQ queue should be consumed.
type ConsumerConfig struct {
	Queue             string
	DeclareDurable    bool
	DeclareAutoDelete bool
	DeclareExclusive  bool
	DeclareNoWait     bool
	DeclareArgs       amqp.Table
	ConsumerName      string
	AutoAck           bool
	Exclusive         bool
	NoLocal           bool
	NoWait            bool
	Args              amqp.Table
}

// PublisherConfig determines how messages should be published to a RabbitMQ queue.
type PublisherConfig struct {
	Queue             string
	DeclareDurable    bool
	DeclareAutoDelete bool
	DeclareExclusive  bool
	DeclareNoWait     bool
	DeclareArgs       amqp.Table
	Exchange          string
	Mandatory         bool
	Immediate         bool
}
