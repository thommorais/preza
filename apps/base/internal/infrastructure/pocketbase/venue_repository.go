package pocketbase

import (
	"context"
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"

	"github.com/pocketbase/pocketbase/core"
)

type VenueRepository struct {
	app core.App
}

func NewVenueRepository(app core.App) *VenueRepository {
	return &VenueRepository{app: app}
}

func (r *VenueRepository) FindByID(_ context.Context, id venue.ID) (*venue.Venue, error) {
	record, err := r.app.FindRecordById("venues", string(id))
	if err != nil {
		return nil, err
	}
	return recordToVenue(record), nil
}

func (r *VenueRepository) FindByPromoterID(_ context.Context, promoterID promoter.ID) ([]*venue.Venue, error) {
	records, err := r.app.FindRecordsByFilter("venues", "promoter_id = {:promoterID}", "-created", 0, 0, map[string]any{"promoterID": string(promoterID)})
	if err != nil {
		return nil, err
	}

	venues := make([]*venue.Venue, len(records))
	for i, rec := range records {
		venues[i] = recordToVenue(rec)
	}
	return venues, nil
}

func (r *VenueRepository) Save(_ context.Context, v *venue.Venue) error {
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
		return err
	}

	v.ID = venue.ID(record.Id)
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
