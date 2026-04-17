package pocketbase_test

import (
	"testing"

	"preza/internal/domain/event"
	"preza/internal/domain/invitation"
	"preza/internal/domain/promoter"
	pb "preza/internal/infrastructure/pocketbase"
)

func TestInvitationRepository_FindByID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	records, err := app.FindRecordsByFilter("invitations", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no invitation records in test data")
	}
	id := records[0].Id

	inv, err := repo.FindByID(bg(), invitation.ID(id))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if string(inv.ID) != id {
		t.Errorf("ID: got %q, want %q", inv.ID, id)
	}
	if inv.Token == "" {
		t.Error("Token should not be empty")
	}
	if inv.EventID == "" {
		t.Error("EventID should not be empty")
	}
	if inv.GuestID == "" {
		t.Error("GuestID should not be empty")
	}
}

func TestInvitationRepository_FindByToken(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	records, err := app.FindRecordsByFilter("invitations", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no invitation records in test data")
	}
	token := records[0].GetString("token")

	inv, err := repo.FindByToken(bg(), token)
	if err != nil {
		t.Fatalf("FindByToken: %v", err)
	}
	if inv.Token != token {
		t.Errorf("Token: got %q, want %q", inv.Token, token)
	}
}

func TestInvitationRepository_FindByToken_Invalid(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	_, err := repo.FindByToken(bg(), "invalidtoken000000000000000000000000")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestInvitationRepository_FindByEventID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	records, err := app.FindRecordsByFilter("invitations", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no invitation records in test data")
	}
	eventID := event.ID(records[0].GetString("event_id"))

	invitations, err := repo.FindByEventID(bg(), eventID)
	if err != nil {
		t.Fatalf("FindByEventID: %v", err)
	}
	if len(invitations) == 0 {
		t.Error("expected at least one invitation for event")
	}
	for _, inv := range invitations {
		if inv.EventID != eventID {
			t.Errorf("EventID mismatch: got %q, want %q", inv.EventID, eventID)
		}
	}
}

func TestInvitationRepository_FindByPromoterID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	records, err := app.FindRecordsByFilter("invitations", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no invitation records in test data")
	}
	promoterID := promoter.ID(records[0].GetString("created_by"))

	invitations, err := repo.FindByPromoterID(bg(), promoterID)
	if err != nil {
		t.Fatalf("FindByPromoterID: %v", err)
	}
	if len(invitations) == 0 {
		t.Error("expected at least one invitation for promoter")
	}
}

func TestInvitationRepository_StatusMapping(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	statuses := []string{"pending", "confirmed", "declined", "used"}
	for _, status := range statuses {
		records, err := app.FindRecordsByFilter(
			"invitations",
			"status = {:s}",
			"+created", 1, 0,
			map[string]any{"s": status},
		)
		if err != nil || len(records) == 0 {
			continue
		}

		inv, err := repo.FindByID(bg(), invitation.ID(records[0].Id))
		if err != nil {
			t.Errorf("FindByID for status %q: %v", status, err)
			continue
		}
		if string(inv.Status) != status {
			t.Errorf("Status: got %q, want %q", inv.Status, status)
		}
	}
}

func TestInvitationRepository_CheckedInAt_Populated(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	records, err := app.FindRecordsByFilter("invitations", "status = 'used'", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Skip("no used invitations in test data")
	}

	inv, err := repo.FindByID(bg(), invitation.ID(records[0].Id))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if inv.CheckedInAt == nil {
		t.Error("CheckedInAt should be set for used invitations")
	}
}

func TestInvitationRepository_Save_CheckIn(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	records, err := app.FindRecordsByFilter("invitations", "status = 'confirmed'", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Skip("no confirmed invitations in test data")
	}

	inv, err := repo.FindByID(bg(), invitation.ID(records[0].Id))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}

	if err := inv.CheckIn(nowUTC()); err != nil {
		t.Fatalf("CheckIn: %v", err)
	}
	if err := repo.Save(bg(), inv); err != nil {
		t.Fatalf("Save: %v", err)
	}

	updated, err := repo.FindByID(bg(), inv.ID)
	if err != nil {
		t.Fatalf("FindByID after save: %v", err)
	}
	if updated.Status != invitation.StatusUsed {
		t.Errorf("Status: got %q, want %q", updated.Status, invitation.StatusUsed)
	}
	if updated.CheckedInAt == nil {
		t.Error("CheckedInAt should be set after check-in")
	}
}

func TestInvitationRepository_Save_AlreadyCheckedIn(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)

	records, err := app.FindRecordsByFilter("invitations", "status = 'used'", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Skip("no used invitations in test data")
	}

	inv, err := repo.FindByID(bg(), invitation.ID(records[0].Id))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}

	if err := inv.CheckIn(nowUTC()); err != invitation.ErrAlreadyCheckedIn {
		t.Errorf("CheckIn on used: got %v, want ErrAlreadyCheckedIn", err)
	}
}
