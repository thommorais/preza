package invitation_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	appinvitation "preza/internal/application/invitation"
	"preza/internal/domain/invitation"
	"preza/internal/domain/logger"
	pb "preza/internal/infrastructure/pocketbase"

	"github.com/pocketbase/pocketbase/tests"
)

type noopLogger struct{}

func (noopLogger) Info(_ string, _ ...logger.Field)  {}
func (noopLogger) Warn(_ string, _ ...logger.Field)  {}
func (noopLogger) Error(_ string, _ ...logger.Field) {}

func newTestApp(t *testing.T) *tests.TestApp {
	t.Helper()

	_, currentFile, _, _ := runtime.Caller(0)
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
	repo := pb.NewInvitationRepository(app, noopLogger{})
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
	repo := pb.NewInvitationRepository(app, noopLogger{})
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
	repo := pb.NewInvitationRepository(app, noopLogger{})
	uc := appinvitation.NewCheckInGuestUseCase(repo)

	_, err := uc.Execute(context.Background(), "invalidtoken000000000000000000000000")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}
