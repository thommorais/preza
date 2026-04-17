package pocketbase

import (
	"context"
	"preza/internal/domain/logger"
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type VenueRepository struct {
	app core.App
	log logger.Logger
}

func NewVenueRepository(app core.App, log logger.Logger) *VenueRepository {
	return &VenueRepository{app: app, log: log}
}

func (r *VenueRepository) FindByID(_ context.Context, id venue.ID) (*venue.Venue, error) {
	start := time.Now()
	record, err := r.app.FindRecordById("venues", string(id))
	if err != nil {
		r.log.Error("venues.FindByID failed", logger.F("id", id), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("venues.FindByID", logger.F("id", id), logger.F("latency_ms", ms(start)))
	return recordToVenue(record), nil
}

func (r *VenueRepository) FindByPromoterID(_ context.Context, promoterID promoter.ID) ([]*venue.Venue, error) {
	start := time.Now()
	records, err := r.app.FindRecordsByFilter("venues", "promoter_id = {:promoterID}", "-created", 0, 0, map[string]any{"promoterID": string(promoterID)})
	if err != nil {
		r.log.Error("venues.FindByPromoterID failed", logger.F("promoter", promoterID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("venues.FindByPromoterID", logger.F("promoter", promoterID), logger.F("count", len(records)), logger.F("latency_ms", ms(start)))

	venues := make([]*venue.Venue, len(records))
	for i, rec := range records {
		venues[i] = recordToVenue(rec)
	}
	return venues, nil
}

func (r *VenueRepository) Save(_ context.Context, v *venue.Venue) error {
	start := time.Now()
	var record *core.Record

	if v.ID != "" {
		existing, err := r.app.FindRecordById("venues", string(v.ID))
		if err != nil {
			return err
		}
		record = existing
	} else {
		col, err := r.app.FindCollectionByNameOrId("venues")
		if err != nil {
			return err
		}
		record = core.NewRecord(col)
	}

	record.Set("promoter_id", string(v.PromoterID))
	record.Set("name", v.Name)
	record.Set("address", v.Address)
	record.Set("city", v.City)
	record.Set("capacity", v.Capacity)

	if err := r.app.Save(record); err != nil {
		r.log.Error("venues.Save failed", logger.F("id", v.ID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return err
	}

	v.ID = venue.ID(record.Id)
	r.log.Info("venues.Save", logger.F("id", v.ID), logger.F("latency_ms", ms(start)))
	return nil
}

func recordToVenue(r *core.Record) *venue.Venue {
	return &venue.Venue{
		ID:         venue.ID(r.Id),
		PromoterID: promoter.ID(r.GetString("promoter_id")),
		Name:       r.GetString("name"),
		Address:    r.GetString("address"),
		City:       r.GetString("city"),
		Capacity:   r.GetInt("capacity"),
		CoverImage: r.GetString("cover_image"),
		CreatedAt:  r.GetDateTime("created").Time(),
		UpdatedAt:  r.GetDateTime("updated").Time(),
	}
}
