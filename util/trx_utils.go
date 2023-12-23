package util

import (
	"github.com/d1vshar/splitgo/dto"
)

func ValidateTrxPayload(payload *dto.TransasctionPayload) bool {
	if payload != nil && len(payload.TransactionParticipants) > 0 {
		var sum float64 = 0
		for _, participant := range payload.TransactionParticipants {
			sum += participant.Amount
		}

		if sum == payload.Transaction.Amount {
			return true
		}
	}

	return false
}
