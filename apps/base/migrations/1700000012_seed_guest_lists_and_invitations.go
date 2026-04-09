package migrations

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		guestListsCol, err := app.FindCollectionByNameOrId("guest_lists")
		if err != nil {
			return err
		}
		invitationsCol, err := app.FindCollectionByNameOrId("invitations")
		if err != nil {
			return err
		}

		promoterRecords, err := app.FindRecordsByFilter("promoters", "id != ''", "+created", 20, 0)
		if err != nil {
			return fmt.Errorf("load promoters: %w", err)
		}
		guestRecords, err := app.FindRecordsByFilter("guests", "id != ''", "+created", 100, 0)
		if err != nil {
			return fmt.Errorf("load guests: %w", err)
		}
		eventRecords, err := app.FindRecordsByFilter("events", "id != ''", "+created", 25, 0)
		if err != nil {
			return fmt.Errorf("load events: %w", err)
		}

		promoterIDs := make([]string, 0, len(promoterRecords))
		for _, r := range promoterRecords {
			promoterIDs = append(promoterIDs, r.Id)
		}
		guestIDs := make([]string, 0, len(guestRecords))
		for _, r := range guestRecords {
			guestIDs = append(guestIDs, r.Id)
		}

		// Separate past and future events by date
		now := time.Now()
		pastEventIDs := make([]string, 0)
		futureEventIDs := make([]string, 0)
		for _, r := range eventRecords {
			dateStr := r.GetString("date")
			t, err := time.Parse("2006-01-02 15:04:05.999Z", dateStr)
			if err != nil {
				t, _ = time.Parse("2006-01-02 15:04:05", dateStr)
			}
			if t.Before(now) {
				pastEventIDs = append(pastEventIDs, r.Id)
			} else {
				futureEventIDs = append(futureEventIDs, r.Id)
			}
		}

		// guest_lists: distribute guests across promoters
		// Each promoter gets ~5 guests on their list
		guestListPairs := make(map[string]bool)
		for i, guestID := range guestIDs {
			promoterID := promoterIDs[i%len(promoterIDs)]
			key := promoterID + ":" + guestID
			if guestListPairs[key] {
				continue
			}
			guestListPairs[key] = true

			glRec := core.NewRecord(guestListsCol)
			glRec.Set("guest_id", guestID)
			glRec.Set("promoter_id", promoterID)
			if err := app.Save(glRec); err != nil {
				continue
			}
		}
		// Add a second promoter for some guests to test shared guests
		for i := range 30 {
			guestID := guestIDs[i%len(guestIDs)]
			promoterID := promoterIDs[(i+5)%len(promoterIDs)]
			key := promoterID + ":" + guestID
			if guestListPairs[key] {
				continue
			}
			guestListPairs[key] = true

			glRec := core.NewRecord(guestListsCol)
			glRec.Set("guest_id", guestID)
			glRec.Set("promoter_id", promoterID)
			if err := app.Save(glRec); err != nil {
				continue
			}
		}

		// invitations for past events: realistic mix of used/confirmed/declined
		// statuses weighted toward "used" since events already happened
		pastStatuses := []string{
			"used", "used", "used", "used",
			"confirmed", "confirmed",
			"declined",
			"pending",
		}
		plusOnesCounts := []int{0, 0, 0, 1, 1, 2, 0, 0}
		plusOneNames := [][]string{
			{},
			{},
			{},
			{"Mariana Alves"},
			{"Carlos Souza"},
			{"Fernanda Lima", "Rafael Costa"},
			{},
			{},
		}

		for ei, eventID := range pastEventIDs {
			// 8-15 guests per past event
			guestCount := 8 + (ei % 8)
			for gi := range guestCount {
				guestIdx := (ei*guestCount + gi) % len(guestIDs)
				promoterIdx := (ei + gi) % len(promoterIDs)
				statusIdx := gi % len(pastStatuses)
				status := pastStatuses[statusIdx]

				token, err := generateToken()
				if err != nil {
					return fmt.Errorf("generate token: %w", err)
				}

				invRec := core.NewRecord(invitationsCol)
				invRec.Set("event_id", eventID)
				invRec.Set("guest_id", guestIDs[guestIdx])
				invRec.Set("created_by", promoterIDs[promoterIdx])
				invRec.Set("status", status)
				invRec.Set("token", token)
				invRec.Set("plus_ones", plusOnesCounts[statusIdx])
				if len(plusOneNames[statusIdx]) > 0 {
					invRec.Set("plus_one_names", plusOneNames[statusIdx])
				}
				if status == "confirmed" || status == "used" {
					rsvpAt := now.AddDate(0, 0, -(ei+1)*10)
					invRec.Set("rsvp_at", rsvpAt.Format("2006-01-02 15:04:05"))
				}
				if status == "used" {
					checkedInAt := now.AddDate(0, 0, -(ei+1)*8)
					invRec.Set("checked_in_at", checkedInAt.Format("2006-01-02 15:04:05"))
				}

				if err := app.Save(invRec); err != nil {
					// skip duplicate token collisions
					continue
				}
			}
		}

		// invitations for future events: pending and confirmed only
		futureStatuses := []string{
			"confirmed", "confirmed", "confirmed",
			"pending", "pending", "pending", "pending",
			"declined",
		}

		for ei, eventID := range futureEventIDs {
			// 10-20 guests per future event
			guestCount := 10 + (ei % 11)
			for gi := range guestCount {
				guestIdx := (ei*guestCount + gi + 50) % len(guestIDs)
				promoterIdx := (ei + gi + 3) % len(promoterIDs)
				statusIdx := gi % len(futureStatuses)
				status := futureStatuses[statusIdx]

				token, err := generateToken()
				if err != nil {
					return fmt.Errorf("generate token: %w", err)
				}

				invRec := core.NewRecord(invitationsCol)
				invRec.Set("event_id", eventID)
				invRec.Set("guest_id", guestIDs[guestIdx])
				invRec.Set("created_by", promoterIDs[promoterIdx])
				invRec.Set("status", status)
				invRec.Set("token", token)
				invRec.Set("plus_ones", plusOnesCounts[statusIdx%len(plusOnesCounts)])
				if status == "confirmed" {
					rsvpAt := now.AddDate(0, 0, -(gi + 1))
					invRec.Set("rsvp_at", rsvpAt.Format("2006-01-02 15:04:05"))
				}

				if err := app.Save(invRec); err != nil {
					continue
				}
			}
		}

		return nil
	}, func(app core.App) error {
		for _, colName := range []string{"invitations", "guest_lists"} {
			records, err := app.FindRecordsByFilter(colName, "id != ''", "-created", 500, 0)
			if err != nil {
				continue
			}
			for _, r := range records {
				_ = app.Delete(r)
			}
		}
		return nil
	})
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
