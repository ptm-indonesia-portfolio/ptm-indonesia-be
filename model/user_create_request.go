package model

type UserCreateRequest struct {
	Name    string  `json:"name" validate:"required,max=100"`
	Email   string  `json:"email" validate:"required,email,max=255"`
	Address *string `json:"address,omitempty"`
	Telp    *string `json:"telp,omitempty" validate:"omitempty,max=30"`
	Status  *int16  `json:"status" validate:"required"`
}
