package model

import "time"

type UserResponse struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	GoogleID  *string    `json:"google_id,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	Address   *string    `json:"address,omitempty"`
	Telp      *string    `json:"telp,omitempty"`
	Status    int16      `json:"status"`
	CreatedBy uint64     `json:"created_by"`
	UpdatedBy uint64     `json:"updated_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
