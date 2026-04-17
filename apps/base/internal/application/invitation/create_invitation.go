package invitation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"preza/internal/domain/event"
	"preza/internal/domain/guest"
	"preza/internal/domain/invitation"
	"preza/internal/domain/promoter"
)

type CreateInvitationInput struct {
	EventID   event.ID
	CreatedBy promoter.ID
	GuestID   guest.ID
	PlusOnes  int
	Notes     string
}

type CreateInvitationUseCase struct {
	invitations invitation.Repository
	events      event.Repository
}

func NewCreateInvitationUseCase(invitations invitation.Repository, events event.Repository) *CreateInvitationUseCase {
	return &CreateInvitationUseCase{invitations: invitations, events: events}
}

func (uc *CreateInvitationUseCase) Execute(ctx context.Context, input CreateInvitationInput) (*invitation.Invitation, error) {
	ev, err := uc.events.FindByID(ctx, input.EventID)
	if err != nil {
		return nil, err
	}

	if ev.Status != event.StatusActive {
		return nil, ErrEventNotActive
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	inv := &invitation.Invitation{
		EventID:   input.EventID,
		CreatedBy: input.CreatedBy,
		GuestID:   input.GuestID,
		Status:    invitation.StatusPending,
		Token:     token,
		PlusOnes:  input.PlusOnes,
		Notes:     input.Notes,
	}

	if err := uc.invitations.Save(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
