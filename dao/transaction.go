package dao

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/d1vshar/splitgo/config"
	"github.com/d1vshar/splitgo/dto"
	"github.com/d1vshar/splitgo/model"
	"github.com/d1vshar/splitgo/util"
	"github.com/google/uuid"
)

type TransactionParticipantJsonAgg struct {
	Amount        float64    `json:"amount"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *string    `json:"created_by"`
	DeletedAt     *time.Time `json:"deleted_at"`
	ID            string     `json:"id"`
	TransactionID string     `json:"transaction_id"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *string    `json:"updated_by"`
	UserID        string     `json:"user_id"`
}

// returns all transactions for user
func (dao AppDao) GetTransactionsForUser(user_id string) (transactions []model.Transaction, err error) {
	// var transaction_participants []model.TransactionParticipant

	rows, err := dao.db.Raw(`
	SELECT
		json_agg (DISTINCT to_jsonb (_trxs_participants.*)) AS transactionParticipants,
		_trxs.*
	FROM
		_trxs
		LEFT JOIN _trxs_participants ON _trxs_participants.transaction_id = _trxs.id
		LEFT JOIN _categories ON _trxs.category_id = _categories.id
	WHERE
		_trxs.id IN (
			SELECT
				_trxs_participants.transaction_id
			FROM
				_trxs_participants
			WHERE
				_trxs_participants.user_id = ?
		)
		AND _trxs_participants.deleted_at IS NULL
		AND _trxs.deleted_at IS NULL
	GROUP BY
		_trxs.id;`, user_id).Rows()

	defer rows.Close()
	for rows.Next() {
		var jsonAggResult string
		var id, created_at, updated_at, deleted_at, category_id sql.NullString
		var transaction model.Transaction

		err := rows.Scan(
			&jsonAggResult,
			&id,
			&created_at,
			&updated_at,
			&deleted_at,
			&transaction.Meta.CreatedBy,
			&transaction.Meta.UpdatedBy,
			&transaction.Description,
			&category_id,
			&transaction.Amount,
		)

		if err != nil {
			return nil, err
		}

		if len(id.String) > 0 {
			transaction.Meta.ID, err = uuid.Parse(id.String)
			if err != nil {
				return nil, err
			}
		}

		if len(category_id.String) > 0 {
			transaction.CategoryID, err = uuid.Parse(category_id.String)
			if err != nil {
				return nil, err
			}
		}

		if len(created_at.String) > 0 {
			created_at_time, err := time.Parse(time.RFC3339, created_at.String)
			if err != nil {
				return nil, err
			}

			transaction.Meta.CreatedAt = created_at_time
		}

		if len(updated_at.String) > 0 {
			updated_at_time, err := time.Parse(time.RFC3339, updated_at.String)
			if err != nil {
				return nil, err
			}

			transaction.Meta.UpdatedAt = updated_at_time
		}

		if len(deleted_at.String) > 0 {
			deleted_at_time, err := time.Parse(time.RFC3339, deleted_at.String)
			if err != nil {
				return nil, err
			}

			transaction.Meta.DeletedAt = &deleted_at_time
		} else {
			transaction.Meta.DeletedAt = nil
		}

		var participants []TransactionParticipantJsonAgg
		// unmarshal json
		if err := json.Unmarshal([]byte(jsonAggResult), &participants); err != nil {
			return nil, err
		}

		trx_participants := convert_json_agg_to_transaction_participants(participants)
		transaction.TransactionParticipants = trx_participants

		transactions = append(transactions, transaction)

	}

	return transactions, err
}

// creates new transaction and participants
func (dao AppDao) CreateNewTransaction(user_id string, payload *dto.TransasctionPayload) (newTrx *model.Transaction, err error) {
	// payload validation
	if !util.ValidateTrxPayload(payload) {
		return nil, errors.New("invalid transaction payload")
	}

	category, err := dao.GetCategory(payload.Transaction.CategoryID)

	if err != nil {
		return nil, errors.New("invalid category id")
	}

	// db transaction
	tx := dao.db.Begin()

	newTrx = &model.Transaction{
		Description: payload.Transaction.Description,
		CategoryID:  payload.Transaction.CategoryID,
		Category:    category,
		Amount:      payload.Transaction.Amount,
		Meta: config.BaseModel{
			CreatedBy: &user_id,
		},
	}

	newTrx.TransactionParticipants = util.Transform(payload.TransactionParticipants, func(p dto.NewTransactionParticipant) model.TransactionParticipant {
		return model.TransactionParticipant{
			UserID: p.UserID,
			Amount: p.Amount,
		}
	})

	res := tx.Preload("Category").Create(&newTrx)

	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	tx.Commit()

	return newTrx, err
}

// updates transaction and participants
func (dao AppDao) UpdateTransaction(trx_id uuid.UUID, user_id string, payload *dto.TransasctionPayload) (o_trx *model.Transaction, err error) {
	// payload validation
	if !util.ValidateTrxPayload(payload) {
		return nil, errors.New("invalid transaction payload")
	}

	category, err := dao.GetCategory(payload.Transaction.CategoryID)

	if err != nil {
		return nil, errors.New("invalid category id")
	}

	// db transaction
	t := dao.db.Begin()

	find_error := t.Preload("Category", "deleted_at IS null").Preload("TransactionParticipants", "deleted_at IS null").Where("deleted_at IS null").First(o_trx, trx_id).Error

	if find_error != nil {
		t.Rollback()
		return nil, errors.New("transaction not found")
	}

	// access validation
	non_participant := !util.Any(o_trx.TransactionParticipants, func(p model.TransactionParticipant) bool {
		return p.UserID == user_id
	})

	if non_participant {
		t.Rollback()
		return nil, errors.New("user is not a authorized")
	}

	o_trx.Description = payload.Transaction.Description
	o_trx.CategoryID = payload.Transaction.CategoryID
	o_trx.Category = category
	o_trx.Amount = payload.Transaction.Amount
	o_trx.Meta.UpdatedBy = &user_id

	// update existing participants and add new ones
	existingParticipants := make(map[string]bool)
	for _, participant := range o_trx.TransactionParticipants {
		existingParticipants[participant.UserID] = true
	}

	for _, newParticipant := range payload.TransactionParticipants {
		if existingParticipants[newParticipant.UserID] {
			// Update existing participant contribution
			for i, participant := range o_trx.TransactionParticipants {
				if participant.UserID == newParticipant.UserID {
					o_trx.TransactionParticipants[i].Amount = newParticipant.Amount
					break
				}
			}
		} else {
			// Add new participant
			o_trx.TransactionParticipants = append(o_trx.TransactionParticipants, model.TransactionParticipant{
				UserID: newParticipant.UserID,
				Amount: newParticipant.Amount,
			})
		}
	}

	res := t.Save(o_trx)

	if res.Error != nil {
		t.Rollback()
		return nil, res.Error
	}

	for _, trxParticipants := range o_trx.TransactionParticipants {
		partRes := t.Save(trxParticipants)

		if partRes.Error != nil {
			t.Rollback()
			return nil, res.Error
		}
	}

	t.Commit()

	return o_trx, err
}

// soft deletes transaction and participants
func (dao AppDao) DeleteTransaction(trx_id uuid.UUID, user_id string) error {
	now := time.Now()

	o_trx := &model.Transaction{}
	t := dao.db.Begin()

	find_error := t.Preload("Category", "deleted_at IS null").Preload("TransactionParticipants", "deleted_at IS null").Where("deleted_at IS null").First(o_trx, trx_id).Error

	if find_error != nil {
		return errors.New("transaction not found")
	}

	// access validation
	non_participant := !util.Any(o_trx.TransactionParticipants, func(p model.TransactionParticipant) bool {
		return p.UserID == user_id
	})

	if non_participant {
		return errors.New("user is not a authorized")
	}

	o_trx.Meta.DeletedAt = &now
	o_trx.Meta.UpdatedBy = &user_id

	res := t.Save(o_trx)

	if res.Error != nil {
		t.Rollback()
		return res.Error
	}

	for _, trxParticipant := range o_trx.TransactionParticipants {
		trxParticipant.Meta.DeletedAt = &now
		trxParticipant.Meta.UpdatedBy = &user_id
		part_res := t.Save(trxParticipant)

		if part_res.Error != nil {
			t.Rollback()
			return res.Error
		}
	}

	t.Commit()

	return nil
}

func convert_json_agg_to_transaction_participants(json_agg []TransactionParticipantJsonAgg) []model.TransactionParticipant {
	var participants []model.TransactionParticipant

	for _, p := range json_agg {
		part_id, err := uuid.Parse(p.ID)
		if err != nil {
			continue
		}

		trx_id, err := uuid.Parse(p.TransactionID)
		if err != nil {
			continue
		}

		participants = append(participants, model.TransactionParticipant{
			Meta: config.BaseModel{
				ID:        part_id,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
				DeletedAt: p.DeletedAt,
				CreatedBy: p.CreatedBy,
				UpdatedBy: p.UpdatedBy,
			},
			TransactionID: trx_id,
			UserID:        p.UserID,
			Amount:        p.Amount,
		})
	}

	return participants
}
