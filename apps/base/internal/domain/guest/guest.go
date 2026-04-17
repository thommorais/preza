package guest

import "time"

type ID string

type Guest struct {
	ID        ID
	Name      string
	Phone     string
	Email     string
	Notes     string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}
