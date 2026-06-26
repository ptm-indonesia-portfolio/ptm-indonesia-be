package model

import "time"

type AuthRefreshToken struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"column:user_id;not null"`
	TokenHash string     `gorm:"column:token_hash;size:255;not null;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedBy uint64     `gorm:"column:created_by;not null;default:0"`
	UpdatedBy uint64     `gorm:"column:updated_by;not null;default:0"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null"`
}

func (AuthRefreshToken) TableName() string {
	return "auth_refresh_tokens"
}
