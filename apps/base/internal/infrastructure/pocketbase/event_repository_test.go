package pocketbase_test

import (
	"testing"
	"time"

	"preza/internal/domain/event"
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"
	pb "preza/internal/infrastructure/pocketbase"
)

func nowUTC() time.Time {
	return time.Now().UTC()
}

func TestEventRepository_FindByID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewEventRepository(app, testLogger())

	records, err := app.FindRecordsByFilter("events", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no event records in test data")
	}
	id := records[0].Id

	ev, err := repo.FindByID(bg(), event.ID(id))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if string(ev.ID) != id {
		t.Errorf("ID: got %q, want %q", ev.ID, id)
	}
	if ev.Name == "" {
		t.Error("Name should not be empty")
	}
	if ev.VenueID == "" {
		t.Error("VenueID should not be empty")
	}
}

func TestEventRepository_FindByVenueID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewEventRepository(app, testLogger())

	records, err := app.FindRecordsByFilter("events", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no event records in test data")
	}
	venueID := venue.ID(records[0].GetString("venue_id"))

	events, err := repo.FindByVenueID(bg(), venueID)
	if err != nil {
		t.Fatalf("FindByVenueID: %v", err)
	}
	if len(events) == 0 {
		t.Error("expected at least one event for venue")
	}
	for _, ev := range events {
		if ev.VenueID != venueID {
			t.Errorf("VenueID mismatch: got %q, want %q", ev.VenueID, venueID)
		}
	}
}

func TestEventRepository_FindByPromoterID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewEventRepository(app, testLogger())

	records, err := app.FindRecordsByFilter("events", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no event records in test data")
	}
	promoterID := promoter.ID(records[0].GetString("created_by"))

	events, err := repo.FindByPromoterID(bg(), promoterID)
	if err != nil {
		t.Fatalf("FindByPromoterID: %v", err)
	}
	if len(events) == 0 {
		t.Error("expected at least one event for promoter")
	}
}

func TestEventRepository_StatusMapping(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewEventRepository(app, testLogger())

	cases := []struct {
		filter string
		want   event.Status
	}{
		{"status = 'active'", event.StatusActive},
		{"status = 'closed'", event.StatusClosed},
		{"status = 'draft'", event.StatusDraft},
	}

	for _, tc := range cases {
		records, err := app.FindRecordsByFilter("events", tc.filter, "+created", 1, 0)
		if err != nil || len(records) == 0 {
			continue
		}

		ev, err := repo.FindByID(bg(), event.ID(records[0].Id))
		if err != nil {
			t.Errorf("FindByID for %q: %v", tc.filter, err)
			continue
		}
		if ev.Status != tc.want {
			t.Errorf("Status: got %q, want %q", ev.Status, tc.want)
		}
	}
}

func TestEventRepository_Save_Update(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewEventRepository(app, testLogger())

	records, err := app.FindRecordsByFilter("events", "status = 'active'", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no active events in test data")
	}

	ev, err := repo.FindByID(bg(), event.ID(records[0].Id))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}

	ev.Status = event.StatusClosed

	if err := repo.Save(bg(), ev); err != nil {
		t.Fatalf("Save: %v", err)
	}

	updated, err := repo.FindByID(bg(), ev.ID)
	if err != nil {
		t.Fatalf("FindByID after save: %v", err)
	}
	if updated.Status != event.StatusClosed {
		t.Errorf("Status: got %q, want %q", updated.Status, event.StatusClosed)
	}
}
