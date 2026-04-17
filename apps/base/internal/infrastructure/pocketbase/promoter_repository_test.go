package pocketbase_test

import (
	"testing"

	"preza/internal/domain/promoter"
	pb "preza/internal/infrastructure/pocketbase"
)

func TestPromoterRepository_FindByID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewPromoterRepository(app, testLogger())

	records, err := app.FindRecordsByFilter("promoters", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no promoter records in test data")
	}
	id := records[0].Id

	p, err := repo.FindByID(bg(), promoter.ID(id))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if string(p.ID) != id {
		t.Errorf("ID: got %q, want %q", p.ID, id)
	}
	if p.Name == "" {
		t.Error("Name should not be empty")
	}
	if p.Type != promoter.TypeVenue && p.Type != promoter.TypeIndependent {
		t.Errorf("unexpected Type: %q", p.Type)
	}
}

func TestPromoterRepository_FindByUserID(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewPromoterRepository(app, testLogger())

	records, err := app.FindRecordsByFilter("promoters", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no promoter records in test data")
	}
	userID := records[0].GetString("user_id")

	p, err := repo.FindByUserID(bg(), userID)
	if err != nil {
		t.Fatalf("FindByUserID: %v", err)
	}
	if p.UserID != userID {
		t.Errorf("UserID: got %q, want %q", p.UserID, userID)
	}
}

func TestPromoterRepository_FindByID_NotFound(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewPromoterRepository(app, testLogger())

	_, err := repo.FindByID(bg(), promoter.ID("doesnotexist0000000"))
	if err == nil {
		t.Error("expected error for nonexistent ID, got nil")
	}
}

func TestPromoterRepository_Save_Update(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewPromoterRepository(app, testLogger())

	records, err := app.FindRecordsByFilter("promoters", "id != ''", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Fatal("no promoter records in test data")
	}
	id := promoter.ID(records[0].Id)

	p, err := repo.FindByID(bg(), id)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}

	original := p.Bio
	p.Bio = "Bio atualizada para teste"

	if err := repo.Save(bg(), p); err != nil {
		t.Fatalf("Save: %v", err)
	}

	updated, err := repo.FindByID(bg(), id)
	if err != nil {
		t.Fatalf("FindByID after save: %v", err)
	}
	if updated.Bio != "Bio atualizada para teste" {
		t.Errorf("Bio: got %q, want %q", updated.Bio, "Bio atualizada para teste")
	}
	_ = original
}
