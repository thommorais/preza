package pocketbase

import (
	"context"
	"preza/internal/domain/guest"
	"preza/internal/domain/logger"
	"preza/internal/domain/promoter"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type GuestRepository struct {
	app core.App
	log logger.Logger
}

func NewGuestRepository(app core.App, log logger.Logger) *GuestRepository {
	return &GuestRepository{app: app, log: log}
}

func (r *GuestRepository) FindByID(_ context.Context, id guest.ID) (*guest.Guest, error) {
	start := time.Now()
	record, err := r.app.FindRecordById("guests", string(id))
	if err != nil {
		r.log.Error("guests.FindByID failed", logger.F("id", id), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("guests.FindByID", logger.F("id", id), logger.F("latency_ms", ms(start)))
	return recordToGuest(record), nil
}

func (r *GuestRepository) FindByPromoterGuestList(_ context.Context, promoterID promoter.ID) ([]*guest.Guest, error) {
	start := time.Now()
	records, err := r.app.FindRecordsByFilter(
		"guests",
		"id in (select guest_id from guest_lists where promoter_id = {:promoterID})",
		"-created", 0, 0,
		map[string]any{"promoterID": string(promoterID)},
	)
	if err != nil {
		r.log.Error("guests.FindByPromoterGuestList failed", logger.F("promoter", promoterID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("guests.FindByPromoterGuestList", logger.F("promoter", promoterID), logger.F("count", len(records)), logger.F("latency_ms", ms(start)))

	guests := make([]*guest.Guest, len(records))
	for i, rec := range records {
		guests[i] = recordToGuest(rec)
	}
	return guests, nil
}

func (r *GuestRepository) Save(_ context.Context, g *guest.Guest) error {
	start := time.Now()
	var record *core.Record

	if g.ID != "" {
		existing, err := r.app.FindRecordById("guests", string(g.ID))
		if err != nil {
			return err
		}
		record = existing
	} else {
		col, err := r.app.FindCollectionByNameOrId("guests")
		if err != nil {
			return err
		}
		record = core.NewRecord(col)
	}

	record.Set("name", g.Name)
	record.Set("phone", g.Phone)
	record.Set("email", g.Email)
	record.Set("notes", g.Notes)
	record.Set("created_by", g.CreatedBy)

	if err := r.app.Save(record); err != nil {
		r.log.Error("guests.Save failed", logger.F("id", g.ID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return err
	}

	g.ID = guest.ID(record.Id)
	r.log.Info("guests.Save", logger.F("id", g.ID), logger.F("latency_ms", ms(start)))
	return nil
}

func recordToGuest(r *core.Record) *guest.Guest {
	return &guest.Guest{
		ID:        guest.ID(r.Id),
		Name:      r.GetString("name"),
		Phone:     r.GetString("phone"),
		Email:     r.GetString("email"),
		Notes:     r.GetString("notes"),
		CreatedBy: r.GetString("created_by"),
		CreatedAt: r.GetDateTime("created").Time(),
		UpdatedAt: r.GetDateTime("updated").Time(),
	}
}
