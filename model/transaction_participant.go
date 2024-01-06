package model

import (
	config "github.com/d1vshar/splitgo/config"
	"github.com/google/uuid"
)

type TransactionParticipant struct {
	Meta          config.BaseModel `gorm:"embedded" json:"meta"`
	TransactionID uuid.UUID        `gorm:"index:trx_part_trx_idx;not null" json:"transactionId"`
	UserID        string           `gorm:"index:trx_part_user_idx;not null" json:"userId"`
	AmountDue     float64          `gorm:"not null" json:"amountDue"`
	AmountPaid    float64          `gorm:"not null" json:"amountPaid"`
}

func (m *TransactionParticipant) TableName() string {
	return "public._trxs_participants"
}
