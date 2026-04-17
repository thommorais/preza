package guest

import (
	"context"
	"preza/internal/domain/promoter"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*Guest, error)
	FindByPromoterGuestList(ctx context.Context, promoterID promoter.ID) ([]*Guest, error)
	Save(ctx context.Context, g *Guest) error
}
