package handler

import (
	"context"
	"log"

	domainmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/domain/models"
	rabbitmqmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

// Handler handles bets received and bets calculated.
type Handler struct {
	betRepository BetRepository
}

// New creates and returns a new Handler.
func New(betRepository BetRepository) *Handler {
	return &Handler{
		betRepository: betRepository,
	}
}

// HandleBetsReceived handles bets received.
func (h *Handler) HandleBetsReceived(
	ctx context.Context,
	betsReceived <-chan rabbitmqmodels.BetReceived,
) <-chan rabbitmqmodels.Bet {
	resultingBets := make(chan rabbitmqmodels.Bet)

	go func() {
		defer close(resultingBets)

		for betReceived := range betsReceived {
			log.Println("Processing bet received, betId:", betReceived.Id)

			// Calculate the domain bet based on the incoming bet received.
			domainBet := domainmodels.Bet{
				Id:                   betReceived.Id,
				CustomerId:           betReceived.CustomerId,
				Status:               "active",
				SelectionId:          betReceived.SelectionId,
				SelectionCoefficient: betReceived.SelectionCoefficient,
				Payment:              betReceived.Payment,
			}

			// Insert the domain bet into the repository.
			err := h.betRepository.InsertBet(ctx, domainBet)
			if err != nil {
				log.Println("Failed to insert bet, error: ", err)
				continue
			}

			// Calculate the resulting bet, which should be published.
			resultingBet := rabbitmqmodels.Bet{
				Id:                   domainBet.Id,
				CustomerId:           domainBet.CustomerId,
				Status:               domainBet.Status,
				SelectionId:          domainBet.SelectionId,
				SelectionCoefficient: domainBet.SelectionCoefficient,
				Payment:              domainBet.Payment,
			}

			select {
			case resultingBets <- resultingBet:
			case <-ctx.Done():
				return
			}
		}
	}()

	return resultingBets
}

// HandleBetsCalculated handles bets calculated.
func (h *Handler) HandleBetsCalculated(
	ctx context.Context,
	betsCalculated <-chan rabbitmqmodels.BetCalculated,
) <-chan rabbitmqmodels.Bet {
	resultingBets := make(chan rabbitmqmodels.Bet)

	go func() {
		defer close(resultingBets)

		for betCalculated := range betsCalculated {
			log.Println("Processing bet calculated, betId:", betCalculated.Id)

			// Fetch the domain bet.
			domainBet, exists, err := h.betRepository.GetBetByID(ctx, betCalculated.Id)
			if err != nil {
				log.Println("Failed to fetch a bet which should be updated, error: ", err)
				continue
			}
			if !exists {
				log.Println("A bet which should be updated does not exist, betId: ", betCalculated.Id)
				continue
			}

			// Update the domain bet based on incoming changes.
			domainBet.Status = betCalculated.Status
			domainBet.Payout = betCalculated.Payout

			// Update the domain bet into the repository.
			err = h.betRepository.UpdateBet(ctx, domainBet)
			if err != nil {
				log.Println("Failed to update bet, error: ", err)
				continue
			}

			// Calculate the resulting bet, which should be published.
			resultingBet := rabbitmqmodels.Bet{
				Id:                   domainBet.Id,
				CustomerId:           domainBet.CustomerId,
				Status:               domainBet.Status,
				SelectionId:          domainBet.SelectionId,
				SelectionCoefficient: domainBet.SelectionCoefficient,
				Payment:              domainBet.Payment,
				Payout:               domainBet.Payout,
			}

			select {
			case resultingBets <- resultingBet:
			case <-ctx.Done():
				return
			}
		}
	}()

	return resultingBets
}
