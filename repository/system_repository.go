package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type SystemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) *SystemRepository {
	return &SystemRepository{
		db: db,
	}
}

func (r *SystemRepository) Ping(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("get sql database handle: %w", err)
	}

	return sqlDB.PingContext(ctx)
}
