package invitation

import (
	"context"
	"preza/internal/domain/event"
	"preza/internal/domain/invitation"
)

type ListByEventUseCase struct {
	invitations invitation.Repository
	events      event.Repository
}

func NewListByEventUseCase(invitations invitation.Repository, events event.Repository) *ListByEventUseCase {
	return &ListByEventUseCase{invitations: invitations, events: events}
}

func (uc *ListByEventUseCase) Execute(ctx context.Context, eventID event.ID) ([]*invitation.Invitation, error) {
	if _, err := uc.events.FindByID(ctx, eventID); err != nil {
		return nil, err
	}
	return uc.invitations.FindByEventID(ctx, eventID)
}
