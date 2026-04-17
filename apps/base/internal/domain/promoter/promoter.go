package promoter

import "time"

type ID string

type Type string

const (
	TypeVenue       Type = "venue"
	TypeIndependent Type = "independent"
)

type Promoter struct {
	ID        ID
	UserID    string
	Name      string
	Bio       string
	Instagram string
	Type      Type
	CreatedAt time.Time
	UpdatedAt time.Time
}
