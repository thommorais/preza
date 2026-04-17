package venue

import (
	"preza/internal/domain/promoter"
	"time"
)

type ID string

type Venue struct {
	ID         ID
	PromoterID promoter.ID
	Name       string
	Address    string
	City       string
	Capacity   int
	CoverImage string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
