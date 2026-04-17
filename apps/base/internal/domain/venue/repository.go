package venue

import (
	"context"
	"preza/internal/domain/promoter"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*Venue, error)
	FindByPromoterID(ctx context.Context, promoterID promoter.ID) ([]*Venue, error)
	Save(ctx context.Context, v *Venue) error
}
