package repository

import (
	"context"
	"time"

	"ptm-indonesia/model"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *model.AuthRefreshToken) error {
	return r.db.WithContext(ctx).Create(refreshToken).Error
}

func (r *RefreshTokenRepository) FindActiveByTokenHash(ctx context.Context, tokenHash string) (*model.AuthRefreshToken, error) {
	var refreshToken model.AuthRefreshToken
	if err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) Rotate(ctx context.Context, currentTokenID uint64, replacement *model.AuthRefreshToken) error {
	now := time.Now()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.AuthRefreshToken{}).
			Where("id = ? AND revoked_at IS NULL AND expires_at > ?", currentTokenID, now).
			Updates(map[string]any{
				"revoked_at": now,
				"updated_by": 0,
			})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return tx.Create(replacement).Error
	})
}

func (r *RefreshTokenRepository) RevokeByTokenHash(ctx context.Context, tokenHash string) error {
	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.AuthRefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Updates(map[string]any{
			"revoked_at": now,
			"updated_by": 0,
		}).Error
}
