package model

type AuthSessionUser struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	GoogleID  *string `json:"google_id,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Address   *string `json:"address,omitempty"`
	Telp      *string `json:"telp,omitempty"`
	Status    int16   `json:"status"`
}
