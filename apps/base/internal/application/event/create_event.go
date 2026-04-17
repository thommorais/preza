package event

import (
	"context"
	"preza/internal/domain/event"
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"
	"time"
)

type CreateEventInput struct {
	VenueID     venue.ID
	CreatedBy   promoter.ID
	Name        string
	Date        time.Time
	MaxCapacity int
	Description string
}

type CreateEventUseCase struct {
	events   event.Repository
	venues   venue.Repository
	promoters promoter.Repository
}

func NewCreateEventUseCase(events event.Repository, venues venue.Repository, promoters promoter.Repository) *CreateEventUseCase {
	return &CreateEventUseCase{events: events, venues: venues, promoters: promoters}
}

func (uc *CreateEventUseCase) Execute(ctx context.Context, input CreateEventInput) (*event.Event, error) {
	if _, err := uc.venues.FindByID(ctx, input.VenueID); err != nil {
		return nil, err
	}

	if _, err := uc.promoters.FindByID(ctx, input.CreatedBy); err != nil {
		return nil, err
	}

	ev := &event.Event{
		VenueID:     input.VenueID,
		CreatedBy:   input.CreatedBy,
		Name:        input.Name,
		Date:        input.Date,
		MaxCapacity: input.MaxCapacity,
		Status:      event.StatusDraft,
		Description: input.Description,
	}

	if err := uc.events.Save(ctx, ev); err != nil {
		return nil, err
	}

	return ev, nil
}
