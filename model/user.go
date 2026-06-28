package model

import "time"

type User struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string     `gorm:"column:name;size:100;not null"`
	Email     string     `gorm:"column:email;size:255;not null"`
	GoogleID  *string    `gorm:"column:google_id;size:255;uniqueIndex"`
	AvatarURL *string    `gorm:"column:avatar_url;type:text"`
	Address   *string    `gorm:"column:address;type:text"`
	Telp      *string    `gorm:"column:telp;size:30"`
	Status    int16      `gorm:"column:status;type:smallint;not null;default:0"`
	StatusRow *int16     `gorm:"column:status_row;type:smallint;default:1"`
	CreatedBy uint64     `gorm:"column:created_by;not null;default:0"`
	UpdatedBy uint64     `gorm:"column:updated_by;not null;default:0"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func ActiveUserStatusRow() *int16 {
	value := int16(1)
	return &value
}
