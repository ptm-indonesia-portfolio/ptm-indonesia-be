package model

const (
	UserStatusNotActive     int16 = 0
	UserStatusSuperAdmin    int16 = 1
	UserStatusFreeMember    int16 = 2
	UserStatusPremiumMember int16 = 3
)

func IsValidUserStatus(status int16) bool {
	switch status {
	case UserStatusNotActive, UserStatusSuperAdmin, UserStatusFreeMember, UserStatusPremiumMember:
		return true
	default:
		return false
	}
}

func IsLoginAllowedUserStatus(status int16) bool {
	return IsValidUserStatus(status) && status != UserStatusNotActive
}
