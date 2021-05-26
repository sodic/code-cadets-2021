package engine

import (
	"context"

	rabbitmqmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

type Handler interface {
	HandleBetsReceived(ctx context.Context, betsReceived <-chan rabbitmqmodels.BetReceived) <-chan rabbitmqmodels.Bet
	HandleBetsCalculated(ctx context.Context, betsCalculated <-chan rabbitmqmodels.BetCalculated) <-chan rabbitmqmodels.Bet
}
