package pocketbase

import (
	"context"
	"preza/internal/domain/promoter"

	"github.com/pocketbase/pocketbase/core"
)

type PromoterRepository struct {
	app core.App
}

func NewPromoterRepository(app core.App) *PromoterRepository {
	return &PromoterRepository{app: app}
}

func (r *PromoterRepository) FindByID(_ context.Context, id promoter.ID) (*promoter.Promoter, error) {
	record, err := r.app.FindRecordById("promoters", string(id))
	if err != nil {
		return nil, err
	}
	return recordToPromoter(record), nil
}

func (r *PromoterRepository) FindByUserID(_ context.Context, userID string) (*promoter.Promoter, error) {
	record, err := r.app.FindFirstRecordByFilter("promoters", "user_id = {:userID}", map[string]any{"userID": userID})
	if err != nil {
		return nil, err
	}
	return recordToPromoter(record), nil
}

func (r *PromoterRepository) Save(_ context.Context, p *promoter.Promoter) error {
	var record *core.Record

	if p.ID != "" {
		existing, err := r.app.FindRecordById("promoters", string(p.ID))
		if err != nil {
			return err
		}
		record = existing
	} else {
		col, err := r.app.FindCollectionByNameOrId("promoters")
		if err != nil {
			return err
		}
		record = core.NewRecord(col)
	}

	record.Set("user_id", p.UserID)
	record.Set("name", p.Name)
	record.Set("bio", p.Bio)
	record.Set("instagram", p.Instagram)
	record.Set("type", string(p.Type))

	if err := r.app.Save(record); err != nil {
		return err
	}

	p.ID = promoter.ID(record.Id)
	return nil
}

func recordToPromoter(r *core.Record) *promoter.Promoter {
	return &promoter.Promoter{
		ID:        promoter.ID(r.Id),
		UserID:    r.GetString("user_id"),
		Name:      r.GetString("name"),
		Bio:       r.GetString("bio"),
		Instagram: r.GetString("instagram"),
		Type:      promoter.Type(r.GetString("type")),
		CreatedAt: r.GetDateTime("created").Time(),
		UpdatedAt: r.GetDateTime("updated").Time(),
	}
}
