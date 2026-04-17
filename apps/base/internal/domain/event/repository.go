package event

import (
	"context"
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*Event, error)
	FindByVenueID(ctx context.Context, venueID venue.ID) ([]*Event, error)
	FindByPromoterID(ctx context.Context, promoterID promoter.ID) ([]*Event, error)
	Save(ctx context.Context, e *Event) error
}
