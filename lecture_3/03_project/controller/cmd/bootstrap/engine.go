package bootstrap

import (
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/cmd/config"
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/domain/mappers"
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/engine"
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/engine/consumer"
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/engine/handler"
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/engine/publisher"
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq"
	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/sqlite"
)

func newBetReceivedConsumer(channel rabbitmq.Channel) *rabbitmq.BetReceivedConsumer {
	betReceivedConsumer, err := rabbitmq.NewBetReceivedConsumer(
		channel,
		rabbitmq.ConsumerConfig{
			Queue:             config.Cfg.Rabbit.ConsumerBetReceivedQueue,
			DeclareDurable:    config.Cfg.Rabbit.DeclareDurable,
			DeclareAutoDelete: config.Cfg.Rabbit.DeclareAutoDelete,
			DeclareExclusive:  config.Cfg.Rabbit.DeclareExclusive,
			DeclareNoWait:     config.Cfg.Rabbit.DeclareNoWait,
			DeclareArgs:       nil,
			ConsumerName:      config.Cfg.Rabbit.ConsumerBetReceivedName,
			AutoAck:           config.Cfg.Rabbit.ConsumerAutoAck,
			Exclusive:         config.Cfg.Rabbit.ConsumerExclusive,
			NoLocal:           config.Cfg.Rabbit.ConsumerNoLocal,
			NoWait:            config.Cfg.Rabbit.ConsumerNoWait,
			Args:              nil,
		},
	)
	if err != nil {
		panic(err)
	}
	return betReceivedConsumer
}

func newBetCalculatedConsumer(channel rabbitmq.Channel) *rabbitmq.BetCalculatedConsumer {
	betCalculatedConsumer, err := rabbitmq.NewBetCalculatedConsumer(
		channel,
		rabbitmq.ConsumerConfig{
			Queue:             config.Cfg.Rabbit.ConsumerBetCalculatedQueue,
			DeclareDurable:    config.Cfg.Rabbit.DeclareDurable,
			DeclareAutoDelete: config.Cfg.Rabbit.DeclareAutoDelete,
			DeclareExclusive:  config.Cfg.Rabbit.DeclareExclusive,
			DeclareNoWait:     config.Cfg.Rabbit.DeclareNoWait,
			DeclareArgs:       nil,
			ConsumerName:      config.Cfg.Rabbit.ConsumerBetCalculatedName,
			AutoAck:           config.Cfg.Rabbit.ConsumerAutoAck,
			Exclusive:         config.Cfg.Rabbit.ConsumerExclusive,
			NoLocal:           config.Cfg.Rabbit.ConsumerNoLocal,
			NoWait:            config.Cfg.Rabbit.ConsumerNoWait,
			Args:              nil,
		},
	)
	if err != nil {
		panic(err)
	}
	return betCalculatedConsumer
}

func newConsumer(betReceivedConsumer consumer.BetReceivedConsumer, betCalculatedConsumer consumer.BetCalculatedConsumer) *consumer.Consumer {
	return consumer.New(betReceivedConsumer, betCalculatedConsumer)
}

func newBetMapper() *mappers.BetMapper {
	return mappers.NewBetMapper()
}

func newBetRepository(dbExecutor sqlite.DatabaseExecutor, betMapper sqlite.BetMapper) *sqlite.BetRepository {
	return sqlite.NewBetRepository(dbExecutor, betMapper)
}

func newHandler(betRepository handler.BetRepository) *handler.Handler {
	return handler.New(betRepository)
}

func newBetPublisher(channel rabbitmq.Channel) *rabbitmq.BetPublisher {
	betPublisher, err := rabbitmq.NewBetPublisher(
		channel,
		rabbitmq.PublisherConfig{
			Queue:             config.Cfg.Rabbit.PublisherBetQueue,
			DeclareDurable:    config.Cfg.Rabbit.DeclareDurable,
			DeclareAutoDelete: config.Cfg.Rabbit.DeclareAutoDelete,
			DeclareExclusive:  config.Cfg.Rabbit.DeclareExclusive,
			DeclareNoWait:     config.Cfg.Rabbit.DeclareNoWait,
			DeclareArgs:       nil,
			Exchange:          config.Cfg.Rabbit.PublisherExchange,
			Mandatory:         config.Cfg.Rabbit.PublisherMandatory,
			Immediate:         config.Cfg.Rabbit.PublisherImmediate,
		},
	)
	if err != nil {
		panic(err)
	}
	return betPublisher
}

func newPublisher(betPublisher publisher.BetPublisher) *publisher.Publisher {
	return publisher.New(betPublisher)
}

func Engine(rabbitMqChannel rabbitmq.Channel, dbExecutor sqlite.DatabaseExecutor) *engine.Engine {
	betReceivedConsumer := newBetReceivedConsumer(rabbitMqChannel)
	betCalculatedConsumer := newBetCalculatedConsumer(rabbitMqChannel)
	consumer := newConsumer(betReceivedConsumer, betCalculatedConsumer)

	betMapper := newBetMapper()
	betRepository := newBetRepository(dbExecutor, betMapper)
	handler := newHandler(betRepository)

	betPublisher := newBetPublisher(rabbitMqChannel)
	publisher := newPublisher(betPublisher)

	return engine.New(consumer, handler, publisher)
}
