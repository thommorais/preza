package event

import (
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"
	"time"
)

type ID string

type Status string

const (
	StatusDraft  Status = "draft"
	StatusActive Status = "active"
	StatusClosed Status = "closed"
)

type Event struct {
	ID          ID
	VenueID     venue.ID
	CreatedBy   promoter.ID
	Name        string
	Date        time.Time
	MaxCapacity int
	Status      Status
	CoverImage  string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
