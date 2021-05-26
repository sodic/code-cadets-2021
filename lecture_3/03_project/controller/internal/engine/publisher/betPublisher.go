package publisher

import (
	"context"

	rabbitmqmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

type BetPublisher interface {
	Publish(ctx context.Context, bets <-chan rabbitmqmodels.Bet)
}
