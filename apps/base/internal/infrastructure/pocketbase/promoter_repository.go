package pocketbase

import (
	"context"
	"preza/internal/domain/logger"
	"preza/internal/domain/promoter"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type PromoterRepository struct {
	app core.App
	log logger.Logger
}

func NewPromoterRepository(app core.App, log logger.Logger) *PromoterRepository {
	return &PromoterRepository{app: app, log: log}
}

func (r *PromoterRepository) FindByID(_ context.Context, id promoter.ID) (*promoter.Promoter, error) {
	start := time.Now()
	record, err := r.app.FindRecordById("promoters", string(id))
	if err != nil {
		r.log.Error("promoters.FindByID failed", logger.F("id", id), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("promoters.FindByID", logger.F("id", id), logger.F("latency_ms", ms(start)))
	return recordToPromoter(record), nil
}

func (r *PromoterRepository) FindByUserID(_ context.Context, userID string) (*promoter.Promoter, error) {
	start := time.Now()
	record, err := r.app.FindFirstRecordByFilter("promoters", "user_id = {:userID}", map[string]any{"userID": userID})
	if err != nil {
		r.log.Error("promoters.FindByUserID failed", logger.F("user_id", userID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("promoters.FindByUserID", logger.F("user_id", userID), logger.F("latency_ms", ms(start)))
	return recordToPromoter(record), nil
}

func (r *PromoterRepository) Save(_ context.Context, p *promoter.Promoter) error {
	start := time.Now()
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
		r.log.Error("promoters.Save failed", logger.F("id", p.ID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return err
	}

	p.ID = promoter.ID(record.Id)
	r.log.Info("promoters.Save", logger.F("id", p.ID), logger.F("latency_ms", ms(start)))
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
