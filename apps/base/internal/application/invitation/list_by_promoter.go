package invitation

import (
	"context"
	"preza/internal/domain/invitation"
	"preza/internal/domain/promoter"
)

type ListByPromoterUseCase struct {
	invitations invitation.Repository
	promoters   promoter.Repository
}

func NewListByPromoterUseCase(invitations invitation.Repository, promoters promoter.Repository) *ListByPromoterUseCase {
	return &ListByPromoterUseCase{invitations: invitations, promoters: promoters}
}

func (uc *ListByPromoterUseCase) Execute(ctx context.Context, promoterID promoter.ID) ([]*invitation.Invitation, error) {
	if _, err := uc.promoters.FindByID(ctx, promoterID); err != nil {
		return nil, err
	}
	return uc.invitations.FindByPromoterID(ctx, promoterID)
}
