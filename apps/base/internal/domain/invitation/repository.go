package invitation

import (
	"context"
	"preza/internal/domain/event"
	"preza/internal/domain/promoter"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*Invitation, error)
	FindByToken(ctx context.Context, token string) (*Invitation, error)
	FindByEventID(ctx context.Context, eventID event.ID) ([]*Invitation, error)
	FindByPromoterID(ctx context.Context, promoterID promoter.ID) ([]*Invitation, error)
	Save(ctx context.Context, inv *Invitation) error
}
