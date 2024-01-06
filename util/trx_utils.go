package util

import (
	"math"
	"time"

	"github.com/d1vshar/splitgo/dto"
	"github.com/d1vshar/splitgo/model"
)

func ValidateTrxPayload(payload *dto.TransasctionPayload) bool {
	if payload != nil && len(payload.TransactionParticipants) > 0 {
		var paid_sum, due_sum float64 = 0, 0
		for _, participant := range payload.TransactionParticipants {
			paid_sum += participant.AmountPaid
			due_sum += participant.AmountDue
		}

		if math.Round(paid_sum) != math.Round(payload.Transaction.Amount) || math.Round(due_sum) != math.Round(payload.Transaction.Amount) {
			return false
		}

		_, err := time.Parse(time.RFC3339, payload.Transaction.Timestamp)
		if err != nil {
			return false
		}
	}

	return true
}

func GenerateSettlementGraph(transactions []model.Transaction) map[string]map[string]float64 {
	var individual_transactions []model.TransactionParticipant

	for _, t := range transactions {
		individual_transactions = append(individual_transactions, t.TransactionParticipants...)
	}

	net_map := constructNetMap(individual_transactions)
	settle_map := make(map[string]map[string]float64)

	settleRecursively(net_map, settle_map)

	return settle_map
}

func constructNetMap(ts []model.TransactionParticipant) map[string]float64 {
	var trxMap map[string]float64 = make(map[string]float64)

	for _, tp := range ts {
		net := tp.AmountPaid - tp.AmountDue
		trxMap[tp.UserID] += net
	}

	return trxMap
}

func settleRecursively(net_map map[string]float64, settle_map map[string]map[string]float64) error {
	if net_map != nil {
		maxDebtor, maxCreditor := findMinAndMax(net_map)

		if net_map[maxDebtor] != 0 && net_map[maxCreditor] != 0 {
			minOf2 := math.Min(-net_map[maxDebtor], net_map[maxCreditor])
			net_map[maxCreditor] -= minOf2
			net_map[maxDebtor] += minOf2

			if settle_map[maxDebtor] == nil {
				settle_map[maxDebtor] = make(map[string]float64)
			}

			settle_map[maxDebtor][maxCreditor] += minOf2
			return settleRecursively(net_map, settle_map)
		}
	}

	return nil
}

func findMinAndMax(net_map map[string]float64) (string, string) {
	var max_key string
	var max_val float64 = 0
	var min_key string
	var min_val float64 = 0

	for k, v := range net_map {
		if v > max_val {
			max_key = k
			max_val = v
		}

		if v < min_val {
			min_key = k
			min_val = v
		}
	}

	return min_key, max_key
}
