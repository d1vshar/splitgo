package model

import (
	config "github.com/d1vshar/splitgo/config"
	"github.com/google/uuid"
)

type TransactionParticipant struct {
	Meta          config.BaseModel `gorm:"embedded" json:"meta"`
	TransactionID uuid.UUID        `gorm:"index:trx_part_trx_idx;not null" json:"transactionId"`
	UserID        string           `gorm:"index:trx_part_user_idx;not null" json:"userId"`
	Amount        float64          `gorm:"non null" json:"amount"`
}

func (m *TransactionParticipant) TableName() string {
	return "public._trxs_participants"
}
