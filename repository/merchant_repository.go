package repository

import (
	"context"
	"errors"
	"github.com/titancodehub/brick-interview/common"
	"github.com/titancodehub/brick-interview/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) *MerchantRepository {
	return &MerchantRepository{
		db: db,
	}
}

func (r *MerchantRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *MerchantRepository) GetByID(ctx context.Context, id string) (model.Merchant, error) {
	var result model.Merchant
	if err := r.db.WithContext(ctx).Where("id=?", id).Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, common.ErrorRecordNotExist
		}
	}
	return result, nil
}

func (r *MerchantRepository) GetForUpdate(ctx context.Context, id string, tx *gorm.DB) (model.Merchant, error) {
	var result model.Merchant
	if err := tx.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("id=?", id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, common.ErrorRecordNotExist
		}
	}

	return result, nil
}

func (r *MerchantRepository) UpdateBalance(ctx context.Context, id string, balance int64, tx *gorm.DB) error {
	if err := tx.Model(model.Merchant{}).WithContext(ctx).Where("id=?", id).Updates(map[string]interface{}{"balance": balance}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrorRecordNotExist
		}
	}
	return nil
}
