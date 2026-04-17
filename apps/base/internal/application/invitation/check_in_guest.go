package invitation

import (
	"context"
	"preza/internal/domain/invitation"
	"time"
)

type CheckInGuestUseCase struct {
	invitations invitation.Repository
}

func NewCheckInGuestUseCase(invitations invitation.Repository) *CheckInGuestUseCase {
	return &CheckInGuestUseCase{invitations: invitations}
}

func (uc *CheckInGuestUseCase) Execute(ctx context.Context, token string) (*invitation.Invitation, error) {
	inv, err := uc.invitations.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := inv.CheckIn(time.Now()); err != nil {
		return nil, err
	}

	if err := uc.invitations.Save(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}
