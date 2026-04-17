package invitation_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	appinvitation "preza/internal/application/invitation"
	"preza/internal/domain/invitation"
	pb "preza/internal/infrastructure/pocketbase"

	"github.com/pocketbase/pocketbase/tests"
)

func newTestApp(t *testing.T) *tests.TestApp {
	t.Helper()

	_, currentFile, _, _ := runtime.Caller(0)
	// currentFile: .../apps/base/internal/application/invitation/...
	// pb_data:     .../apps/base/pb_data
	pbDataDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "pb_data")

	app, err := tests.NewTestApp(pbDataDir)
	if err != nil {
		t.Fatalf("create test app: %v", err)
	}

	t.Cleanup(app.Cleanup)
	return app
}

func TestCheckInGuest_Success(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)
	uc := appinvitation.NewCheckInGuestUseCase(repo)

	records, err := app.FindRecordsByFilter("invitations", "status = 'confirmed'", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Skip("no confirmed invitations in test data")
	}
	token := records[0].GetString("token")

	inv, err := uc.Execute(context.Background(), token)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if inv.Status != invitation.StatusUsed {
		t.Errorf("Status: got %q, want %q", inv.Status, invitation.StatusUsed)
	}
	if inv.CheckedInAt == nil {
		t.Error("CheckedInAt should be set after check-in")
	}
}

func TestCheckInGuest_AlreadyCheckedIn(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)
	uc := appinvitation.NewCheckInGuestUseCase(repo)

	records, err := app.FindRecordsByFilter("invitations", "status = 'used'", "+created", 1, 0)
	if err != nil || len(records) == 0 {
		t.Skip("no used invitations in test data")
	}
	token := records[0].GetString("token")

	_, err = uc.Execute(context.Background(), token)
	if err != invitation.ErrAlreadyCheckedIn {
		t.Errorf("expected ErrAlreadyCheckedIn, got %v", err)
	}
}

func TestCheckInGuest_InvalidToken(t *testing.T) {
	app := newTestApp(t)
	repo := pb.NewInvitationRepository(app)
	uc := appinvitation.NewCheckInGuestUseCase(repo)

	_, err := uc.Execute(context.Background(), "invalidtoken000000000000000000000000")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}
