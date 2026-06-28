package services

import (
	"errors"
	"strings"

	"github.com/lib/pq"
)

const userEmailUniqueConstraintName = "idx_users_email_status_row_unique"

func mapUserPersistenceError(err error) error {
	if isUserEmailUniqueConstraintError(err) {
		return ErrUserEmailAlreadyUsed
	}

	return err
}

func isUserEmailUniqueConstraintError(err error) bool {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return false
	}

	if pqErr.Code != "23505" {
		return false
	}

	return strings.EqualFold(pqErr.Constraint, userEmailUniqueConstraintName) ||
		strings.EqualFold(pqErr.Constraint, "idx_users_email_active_unique") ||
		strings.EqualFold(pqErr.Constraint, "users_email_key")
}
