package models

// BetCalculated represents a DTO for calculated bets.
type BetCalculated struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Payout float64 `json:"payout"`
}
