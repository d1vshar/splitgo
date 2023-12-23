package model

import (
	config "github.com/d1vshar/splitgo/config"
	"github.com/google/uuid"
)

type Transaction struct {
	Meta                    config.BaseModel         `gorm:"embedded"`
	Description             string                   `gorm:"non null" json:"description"`
	CategoryID              uuid.UUID                `gorm:"index:trx_category_idx;non null" json:"categoryId"`
	Category                Category                 `json:"-"`
	TransactionParticipants []TransactionParticipant `json:"transactionParticipants"`
	Amount                  float64                  `gorm:"non null" json:"amount"`
}

func (m *Transaction) TableName() string {
	return "public._trxs"
}
