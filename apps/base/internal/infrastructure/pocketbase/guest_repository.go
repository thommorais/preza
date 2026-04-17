package pocketbase

import (
	"context"
	"preza/internal/domain/guest"
	"preza/internal/domain/promoter"

	"github.com/pocketbase/pocketbase/core"
)

type GuestRepository struct {
	app core.App
}

func NewGuestRepository(app core.App) *GuestRepository {
	return &GuestRepository{app: app}
}

func (r *GuestRepository) FindByID(_ context.Context, id guest.ID) (*guest.Guest, error) {
	record, err := r.app.FindRecordById("guests", string(id))
	if err != nil {
		return nil, err
	}
	return recordToGuest(record), nil
}

func (r *GuestRepository) FindByPromoterGuestList(_ context.Context, promoterID promoter.ID) ([]*guest.Guest, error) {
	records, err := r.app.FindRecordsByFilter(
		"guests",
		"id in (select guest_id from guest_lists where promoter_id = {:promoterID})",
		"-created", 0, 0,
		map[string]any{"promoterID": string(promoterID)},
	)
	if err != nil {
		return nil, err
	}

	guests := make([]*guest.Guest, len(records))
	for i, rec := range records {
		guests[i] = recordToGuest(rec)
	}
	return guests, nil
}

func (r *GuestRepository) Save(_ context.Context, g *guest.Guest) error {
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
		return err
	}

	g.ID = guest.ID(record.Id)
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
