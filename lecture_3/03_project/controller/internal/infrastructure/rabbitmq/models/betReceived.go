package models

// BetReceived represents a DTO for received bets.
type BetReceived struct {
	Id                   string  `json:"id"`
	CustomerId           string  `json:"customerId"`
	SelectionId          string  `json:"selectionId"`
	SelectionCoefficient float64 `json:"selectionCoefficient"`
	Payment              float64 `json:"payment"`
}
