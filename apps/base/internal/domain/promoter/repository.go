package promoter

import "context"

type Repository interface {
	FindByID(ctx context.Context, id ID) (*Promoter, error)
	FindByUserID(ctx context.Context, userID string) (*Promoter, error)
	Save(ctx context.Context, p *Promoter) error
}
