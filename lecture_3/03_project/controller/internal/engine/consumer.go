package engine

import (
	"context"

	rabbitmqmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

type Consumer interface {
	ConsumeBetsReceived(ctx context.Context) (<-chan rabbitmqmodels.BetReceived, error)
	ConsumeBetsCalculated(ctx context.Context) (<-chan rabbitmqmodels.BetCalculated, error)
}
