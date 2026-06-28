package model

type UserListResponse struct {
	Items []UserResponse `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}
