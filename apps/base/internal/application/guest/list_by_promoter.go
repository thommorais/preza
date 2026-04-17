package guest

import (
	"context"
	"preza/internal/domain/guest"
	"preza/internal/domain/promoter"
)

type ListByPromoterUseCase struct {
	guests    guest.Repository
	promoters promoter.Repository
}

func NewListByPromoterUseCase(guests guest.Repository, promoters promoter.Repository) *ListByPromoterUseCase {
	return &ListByPromoterUseCase{guests: guests, promoters: promoters}
}

func (uc *ListByPromoterUseCase) Execute(ctx context.Context, promoterID promoter.ID) ([]*guest.Guest, error) {
	if _, err := uc.promoters.FindByID(ctx, promoterID); err != nil {
		return nil, err
	}
	return uc.guests.FindByPromoterGuestList(ctx, promoterID)
}
