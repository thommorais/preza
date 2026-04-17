package invitation

import "errors"

var (
	ErrAlreadyCheckedIn = errors.New("invitation already checked in")
	ErrNotFound         = errors.New("invitation not found")
	ErrInvalidToken     = errors.New("invalid invitation token")
)
