package model

import "time"

const (
	UserStatusNotActive     int16 = 0
	UserStatusSuperAdmin    int16 = 1
	UserStatusFreeMember    int16 = 2
	UserStatusPremiumMember int16 = 3
)

type User struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string     `gorm:"column:name;size:100;not null"`
	Email     string     `gorm:"column:email;size:255;not null;uniqueIndex"`
	GoogleID  *string    `gorm:"column:google_id;size:255;uniqueIndex"`
	AvatarURL *string    `gorm:"column:avatar_url;type:text"`
	Address   *string    `gorm:"column:address;type:text"`
	Telp      *string    `gorm:"column:telp;size:30"`
	Status    int16      `gorm:"column:status;type:smallint;not null;default:0"`
	CreatedBy uint64     `gorm:"column:created_by;not null;default:0"`
	UpdatedBy uint64     `gorm:"column:updated_by;not null;default:0"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
