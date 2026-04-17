package pocketbase

import (
	"context"
	"preza/internal/domain/event"
	"preza/internal/domain/guest"
	"preza/internal/domain/invitation"
	"preza/internal/domain/logger"
	"preza/internal/domain/promoter"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type InvitationRepository struct {
	app core.App
	log logger.Logger
}

func NewInvitationRepository(app core.App, log logger.Logger) *InvitationRepository {
	return &InvitationRepository{app: app, log: log}
}

func (r *InvitationRepository) FindByID(_ context.Context, id invitation.ID) (*invitation.Invitation, error) {
	start := time.Now()
	record, err := r.app.FindRecordById("invitations", string(id))
	if err != nil {
		r.log.Error("invitations.FindByID failed", logger.F("id", id), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("invitations.FindByID", logger.F("id", id), logger.F("latency_ms", ms(start)))
	return recordToInvitation(record), nil
}

func (r *InvitationRepository) FindByToken(_ context.Context, token string) (*invitation.Invitation, error) {
	start := time.Now()
	record, err := r.app.FindFirstRecordByFilter("invitations", "token = {:token}", map[string]any{"token": token})
	if err != nil {
		r.log.Error("invitations.FindByToken failed", logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("invitations.FindByToken", logger.F("id", record.Id), logger.F("latency_ms", ms(start)))
	return recordToInvitation(record), nil
}

func (r *InvitationRepository) FindByEventID(_ context.Context, eventID event.ID) ([]*invitation.Invitation, error) {
	start := time.Now()
	records, err := r.app.FindRecordsByFilter("invitations", "event_id = {:eventID}", "-created", 0, 0, map[string]any{"eventID": string(eventID)})
	if err != nil {
		r.log.Error("invitations.FindByEventID failed", logger.F("event", eventID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("invitations.FindByEventID", logger.F("event", eventID), logger.F("count", len(records)), logger.F("latency_ms", ms(start)))

	invitations := make([]*invitation.Invitation, len(records))
	for i, rec := range records {
		invitations[i] = recordToInvitation(rec)
	}
	return invitations, nil
}

func (r *InvitationRepository) FindByPromoterID(_ context.Context, promoterID promoter.ID) ([]*invitation.Invitation, error) {
	start := time.Now()
	records, err := r.app.FindRecordsByFilter("invitations", "created_by = {:promoterID}", "-created", 0, 0, map[string]any{"promoterID": string(promoterID)})
	if err != nil {
		r.log.Error("invitations.FindByPromoterID failed", logger.F("promoter", promoterID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return nil, err
	}
	r.log.Info("invitations.FindByPromoterID", logger.F("promoter", promoterID), logger.F("count", len(records)), logger.F("latency_ms", ms(start)))

	invitations := make([]*invitation.Invitation, len(records))
	for i, rec := range records {
		invitations[i] = recordToInvitation(rec)
	}
	return invitations, nil
}

func (r *InvitationRepository) Save(_ context.Context, inv *invitation.Invitation) error {
	start := time.Now()
	var record *core.Record

	if inv.ID != "" {
		existing, err := r.app.FindRecordById("invitations", string(inv.ID))
		if err != nil {
			return err
		}
		record = existing
	} else {
		col, err := r.app.FindCollectionByNameOrId("invitations")
		if err != nil {
			return err
		}
		record = core.NewRecord(col)
	}

	record.Set("event_id", string(inv.EventID))
	record.Set("created_by", string(inv.CreatedBy))
	record.Set("guest_id", string(inv.GuestID))
	record.Set("status", string(inv.Status))
	record.Set("token", inv.Token)
	record.Set("plus_ones", inv.PlusOnes)
	record.Set("plus_one_names", inv.PlusOneNames)
	record.Set("notes", inv.Notes)

	if inv.RSVPAt != nil {
		record.Set("rsvp_at", inv.RSVPAt)
	}
	if inv.CheckedInAt != nil {
		record.Set("checked_in_at", inv.CheckedInAt)
	}

	if err := r.app.Save(record); err != nil {
		r.log.Error("invitations.Save failed", logger.F("id", inv.ID), logger.F("latency_ms", ms(start)), logger.F("error", err))
		return err
	}

	inv.ID = invitation.ID(record.Id)
	r.log.Info("invitations.Save", logger.F("id", inv.ID), logger.F("latency_ms", ms(start)))
	return nil
}

func recordToInvitation(r *core.Record) *invitation.Invitation {
	inv := &invitation.Invitation{
		ID:           invitation.ID(r.Id),
		EventID:      event.ID(r.GetString("event_id")),
		CreatedBy:    promoter.ID(r.GetString("created_by")),
		GuestID:      guest.ID(r.GetString("guest_id")),
		Status:       invitation.Status(r.GetString("status")),
		Token:        r.GetString("token"),
		PlusOnes:     r.GetInt("plus_ones"),
		PlusOneNames: r.GetStringSlice("plus_one_names"),
		Notes:        r.GetString("notes"),
		CreatedAt:    r.GetDateTime("created").Time(),
		UpdatedAt:    r.GetDateTime("updated").Time(),
	}

	if rsvpAt := r.GetDateTime("rsvp_at"); !rsvpAt.IsZero() {
		t := rsvpAt.Time()
		inv.RSVPAt = &t
	}
	if checkedInAt := r.GetDateTime("checked_in_at"); !checkedInAt.IsZero() {
		t := checkedInAt.Time()
		inv.CheckedInAt = &t
	}

	return inv
}
