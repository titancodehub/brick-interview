package repository

import (
	"context"
	"github.com/titancodehub/brick-interview/model"
	"gorm.io/gorm"
)

type EntriesRepository struct {
	db *gorm.DB
}

func NewEntriesRepository(db *gorm.DB) *EntriesRepository {
	return &EntriesRepository{
		db: db,
	}
}

func (r *EntriesRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *EntriesRepository) Create(ctx context.Context, entry *model.LedgerEntry, tx *gorm.DB) error {
	executor := r.db
	if tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(entry).Error; err != nil {
		return err
	}

	return nil
}
