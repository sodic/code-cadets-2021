package consumer

import (
	"context"

	rabbitmqmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

type BetCalculatedConsumer interface {
	Consume(ctx context.Context) (<-chan rabbitmqmodels.BetCalculated, error)
}
