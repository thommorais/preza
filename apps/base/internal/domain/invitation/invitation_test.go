package invitation_test

import (
	"testing"
	"time"

	"preza/internal/domain/invitation"
)

func pendingInvitation() *invitation.Invitation {
	return &invitation.Invitation{
		ID:        "inv001",
		Status:    invitation.StatusPending,
		Token:     "abc123",
		PlusOnes:  0,
	}
}

func TestCheckIn_PendingInvitation(t *testing.T) {
	inv := pendingInvitation()
	now := time.Now()

	if err := inv.CheckIn(now); err != nil {
		t.Fatalf("CheckIn on pending: %v", err)
	}
	if inv.Status != invitation.StatusUsed {
		t.Errorf("Status: got %q, want %q", inv.Status, invitation.StatusUsed)
	}
	if inv.CheckedInAt == nil {
		t.Error("CheckedInAt should be set")
	}
	if !inv.CheckedInAt.Equal(now) {
		t.Errorf("CheckedInAt: got %v, want %v", inv.CheckedInAt, now)
	}
}

func TestCheckIn_ConfirmedInvitation(t *testing.T) {
	inv := pendingInvitation()
	inv.Status = invitation.StatusConfirmed

	if err := inv.CheckIn(time.Now()); err != nil {
		t.Fatalf("CheckIn on confirmed: %v", err)
	}
	if inv.Status != invitation.StatusUsed {
		t.Errorf("Status: got %q, want %q", inv.Status, invitation.StatusUsed)
	}
}

func TestCheckIn_AlreadyCheckedIn(t *testing.T) {
	inv := pendingInvitation()
	inv.Status = invitation.StatusUsed

	err := inv.CheckIn(time.Now())
	if err != invitation.ErrAlreadyCheckedIn {
		t.Errorf("got %v, want ErrAlreadyCheckedIn", err)
	}
}

func TestCheckIn_IdempotentStatus(t *testing.T) {
	inv := pendingInvitation()
	inv.Status = invitation.StatusUsed
	originalStatus := inv.Status

	_ = inv.CheckIn(time.Now())

	if inv.Status != originalStatus {
		t.Errorf("Status changed after rejected check-in: got %q", inv.Status)
	}
}
