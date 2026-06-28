package model

const (
	UserListSortByName      = "name"
	UserListSortByEmail     = "email"
	UserListSortByStatus    = "status"
	UserListSortByCreatedAt = "created_at"
	UserListSortByUpdatedAt = "updated_at"

	UserListSortDirectionAsc  = "asc"
	UserListSortDirectionDesc = "desc"
)

func IsValidUserListSortBy(value string) bool {
	switch value {
	case UserListSortByName, UserListSortByEmail, UserListSortByStatus, UserListSortByCreatedAt, UserListSortByUpdatedAt:
		return true
	default:
		return false
	}
}

func IsValidUserListSortDirection(value string) bool {
	switch value {
	case UserListSortDirectionAsc, UserListSortDirectionDesc:
		return true
	default:
		return false
	}
}
