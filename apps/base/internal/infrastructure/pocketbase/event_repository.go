package pocketbase

import (
	"context"
	"preza/internal/domain/event"
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"

	"github.com/pocketbase/pocketbase/core"
)

type EventRepository struct {
	app core.App
}

func NewEventRepository(app core.App) *EventRepository {
	return &EventRepository{app: app}
}

func (r *EventRepository) FindByID(_ context.Context, id event.ID) (*event.Event, error) {
	record, err := r.app.FindRecordById("events", string(id))
	if err != nil {
		return nil, err
	}
	return recordToEvent(record), nil
}

func (r *EventRepository) FindByVenueID(_ context.Context, venueID venue.ID) ([]*event.Event, error) {
	records, err := r.app.FindRecordsByFilter("events", "venue_id = {:venueID}", "-date", 0, 0, map[string]any{"venueID": string(venueID)})
	if err != nil {
		return nil, err
	}

	events := make([]*event.Event, len(records))
	for i, rec := range records {
		events[i] = recordToEvent(rec)
	}
	return events, nil
}

func (r *EventRepository) FindByPromoterID(_ context.Context, promoterID promoter.ID) ([]*event.Event, error) {
	records, err := r.app.FindRecordsByFilter("events", "created_by = {:promoterID}", "-date", 0, 0, map[string]any{"promoterID": string(promoterID)})
	if err != nil {
		return nil, err
	}

	events := make([]*event.Event, len(records))
	for i, rec := range records {
		events[i] = recordToEvent(rec)
	}
	return events, nil
}

func (r *EventRepository) Save(_ context.Context, e *event.Event) error {
	var record *core.Record

	if e.ID != "" {
		existing, err := r.app.FindRecordById("events", string(e.ID))
		if err != nil {
			return err
		}
		record = existing
	} else {
		col, err := r.app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}
		record = core.NewRecord(col)
	}

	record.Set("venue_id", string(e.VenueID))
	record.Set("created_by", string(e.CreatedBy))
	record.Set("name", e.Name)
	record.Set("date", e.Date)
	record.Set("max_capacity", e.MaxCapacity)
	record.Set("status", string(e.Status))
	record.Set("description", e.Description)

	if err := r.app.Save(record); err != nil {
		return err
	}

	e.ID = event.ID(record.Id)
	return nil
}

func recordToEvent(r *core.Record) *event.Event {
	return &event.Event{
		ID:          event.ID(r.Id),
		VenueID:     venue.ID(r.GetString("venue_id")),
		CreatedBy:   promoter.ID(r.GetString("created_by")),
		Name:        r.GetString("name"),
		Date:        r.GetDateTime("date").Time(),
		MaxCapacity: r.GetInt("max_capacity"),
		Status:      event.Status(r.GetString("status")),
		CoverImage:  r.GetString("cover_image"),
		Description: r.GetString("description"),
		CreatedAt:   r.GetDateTime("created").Time(),
		UpdatedAt:   r.GetDateTime("updated").Time(),
	}
}
