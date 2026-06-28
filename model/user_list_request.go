package model

import "strings"

const (
	DefaultUserListPage  = 1
	DefaultUserListLimit = 10
	MaxUserListLimit     = 100
)

type UserListRequest struct {
	Page          int    `query:"page" validate:"omitempty,min=1"`
	Limit         int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search        string `query:"search" validate:"omitempty,max=255"`
	SortBy        string `query:"sort_by" validate:"omitempty,max=50"`
	SortDirection string `query:"sort_direction" validate:"omitempty,max=4"`
}

func (r *UserListRequest) Normalize() {
	if r.Page <= 0 {
		r.Page = DefaultUserListPage
	}

	if r.Limit <= 0 {
		r.Limit = DefaultUserListLimit
	}

	r.Search = normalizeSearchKeyword(r.Search)
	r.SortBy = strings.ToLower(strings.TrimSpace(r.SortBy))
	r.SortDirection = strings.ToLower(strings.TrimSpace(r.SortDirection))
}

func (r UserListRequest) Offset() int {
	return (r.Page - 1) * r.Limit
}

func normalizeSearchKeyword(value string) string {
	return strings.TrimSpace(value)
}
