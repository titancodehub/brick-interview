package repository

import (
	"context"
	"errors"
	"github.com/titancodehub/brick-interview/common"
	"github.com/titancodehub/brick-interview/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *TransactionRepository) Create(ctx context.Context, txn *model.Transaction, tx *gorm.DB) error {
	executor := r.db
	if tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(txn).Error; err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) GetByReference(ctx context.Context, merchantID string, reference string) (model.Transaction, error) {
	var result model.Transaction
	if err := r.db.WithContext(ctx).Where("merchant_id=? AND reference=?", merchantID, reference).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, common.ErrorRecordNotExist
		}
		return result, err
	}

	return result, nil
}

func (r *TransactionRepository) GetByBankReference(ctx context.Context, reference string) (model.Transaction, error) {
	var result model.Transaction
	if err := r.db.WithContext(ctx).Where("bank_reference=?", reference).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, common.ErrorRecordNotExist
		}
		return result, err
	}

	return result, nil
}

func (r *TransactionRepository) GetForUpdate(ctx context.Context, id string, tx *gorm.DB) (model.Transaction, error) {
	var result model.Transaction
	if err := tx.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("id=?", id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, common.ErrorRecordNotExist
		}
	}

	return result, nil
}

func (r *TransactionRepository) Update(ctx context.Context, txn *model.Transaction, tx *gorm.DB) error {
	executor := r.db
	if tx != nil {
		executor = tx
	}

	timestamp := time.Now()
	if err := executor.WithContext(ctx).Model(txn).Updates(map[string]interface{}{"status": txn.Status, "updated": timestamp, "bank_reference": txn.BankReference}).Error; err != nil {
		return err
	}

	return nil
}
