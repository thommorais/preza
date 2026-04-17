package event

import (
	"context"
	"preza/internal/domain/event"
	"preza/internal/domain/promoter"
)

type ListByPromoterUseCase struct {
	events    event.Repository
	promoters promoter.Repository
}

func NewListByPromoterUseCase(events event.Repository, promoters promoter.Repository) *ListByPromoterUseCase {
	return &ListByPromoterUseCase{events: events, promoters: promoters}
}

func (uc *ListByPromoterUseCase) Execute(ctx context.Context, promoterID promoter.ID) ([]*event.Event, error) {
	if _, err := uc.promoters.FindByID(ctx, promoterID); err != nil {
		return nil, err
	}
	return uc.events.FindByPromoterID(ctx, promoterID)
}
