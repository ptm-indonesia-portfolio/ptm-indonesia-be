package model

import "time"

type SystemSetting struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Key       string    `gorm:"column:key;size:100;not null;uniqueIndex"`
	Value     string    `gorm:"column:value;type:text;not null"`
	CreatedBy uint64    `gorm:"column:created_by;not null;default:0"`
	UpdatedBy uint64    `gorm:"column:updated_by;not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}
