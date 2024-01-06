package model

import (
	"time"

	config "github.com/d1vshar/splitgo/config"
	"github.com/google/uuid"
)

type Transaction struct {
	Meta                    config.BaseModel         `gorm:"embedded"`
	Description             string                   `gorm:"non null" json:"description"`
	CategoryID              uuid.UUID                `gorm:"index:trx_category_idx;non null" json:"categoryId"`
	Category                Category                 `json:"-"`
	Timestamp               time.Time                `json:"timestamp"`
	TransactionParticipants []TransactionParticipant `json:"transactionParticipants"`
	Amount                  float64                  `gorm:"not null" json:"amount"`
}

func (m *Transaction) TableName() string {
	return "public._trxs"
}
