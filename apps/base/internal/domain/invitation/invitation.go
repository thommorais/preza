package invitation

import (
	"preza/internal/domain/event"
	"preza/internal/domain/guest"
	"preza/internal/domain/promoter"
	"time"
)

type ID string

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusDeclined  Status = "declined"
	StatusUsed      Status = "used"
)

type Invitation struct {
	ID           ID
	EventID      event.ID
	CreatedBy    promoter.ID
	GuestID      guest.ID
	Status       Status
	Token        string
	PlusOnes     int
	PlusOneNames []string
	RSVPAt       *time.Time
	CheckedInAt  *time.Time
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (i *Invitation) CheckIn(at time.Time) error {
	if i.Status == StatusUsed {
		return ErrAlreadyCheckedIn
	}
	i.Status = StatusUsed
	i.CheckedInAt = &at
	return nil
}
