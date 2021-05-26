package engine

import (
	"context"
	"log"
)

// Engine is the main component, responsible for consuming received bets and calculated bets,
// processing them and publishing the resulting bets.
type Engine struct {
	consumer  Consumer
	handler   Handler
	publisher Publisher
}

// New creates and returns a new engine.
func New(consumer Consumer, handler Handler, publisher Publisher) *Engine {
	return &Engine{
		consumer:  consumer,
		handler:   handler,
		publisher: publisher,
	}
}

// Start will run the engine.
func (e *Engine) Start(ctx context.Context) {
	err := e.processBetsReceived(ctx)
	if err != nil {
		log.Println("Engine failed to process bets received:", err)
		return
	}

	err = e.processBetsCalculated(ctx)
	if err != nil {
		log.Println("Engine failed to process bets calculated:", err)
		return
	}

	<-ctx.Done()
}

func (e *Engine) processBetsReceived(ctx context.Context) error {
	consumedBetsReceived, err := e.consumer.ConsumeBetsReceived(ctx)
	if err != nil {
		return err
	}

	resultingBets := e.handler.HandleBetsReceived(ctx, consumedBetsReceived)
	e.publisher.PublishBets(ctx, resultingBets)

	return nil
}

func (e *Engine) processBetsCalculated(ctx context.Context) error {
	consumedBetsCalculated, err := e.consumer.ConsumeBetsCalculated(ctx)
	if err != nil {
		return err
	}

	resultingBets := e.handler.HandleBetsCalculated(ctx, consumedBetsCalculated)
	e.publisher.PublishBets(ctx, resultingBets)

	return nil
}
