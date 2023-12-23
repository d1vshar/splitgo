package dto

import "github.com/google/uuid"

type TransasctionPayload struct {
	Transaction             NewTransaction              `json:"transaction"`
	TransactionParticipants []NewTransactionParticipant `json:"transactionParticipants"`
}

type NewTransaction struct {
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"categoryId"`
	Amount      float64   `json:"amount"`
}

type NewTransactionParticipant struct {
	UserID string  `json:"userId"`
	Amount float64 `json:"amount"`
}
