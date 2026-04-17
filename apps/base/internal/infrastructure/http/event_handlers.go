package http

import (
	"encoding/json"
	appevent "preza/internal/application/event"
	"preza/internal/domain/logger"
	"preza/internal/domain/promoter"
	"preza/internal/domain/venue"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

func (h *Hooks) handleCreateEvent(e *core.RequestEvent) error {
	requestID := e.Request.Header.Get("X-Request-Id")

	var body struct {
		VenueID     string `json:"venue_id"`
		CreatedBy   string `json:"created_by"`
		Name        string `json:"name"`
		Date        string `json:"date"`
		MaxCapacity int    `json:"max_capacity"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(e.Request.Body).Decode(&body); err != nil {
		h.log.Warn("invalid create event body",
			logger.F("request_id", requestID),
			logger.F("kind", logger.KindValidation),
			logger.F("error", err),
		)
		return e.BadRequestError("invalid request body", err)
	}

	date, err := time.Parse(time.RFC3339, body.Date)
	if err != nil {
		h.log.Warn("invalid event date",
			logger.F("request_id", requestID),
			logger.F("date", body.Date),
			logger.F("kind", logger.KindValidation),
			logger.F("error", err),
		)
		return e.BadRequestError("invalid date format, use RFC3339", err)
	}

	ev, err := h.createEvent.Execute(e.Request.Context(), appevent.CreateEventInput{
		VenueID:     venue.ID(body.VenueID),
		CreatedBy:   promoter.ID(body.CreatedBy),
		Name:        body.Name,
		Date:        date,
		MaxCapacity: body.MaxCapacity,
		Description: body.Description,
	})
	if err != nil {
		h.log.Error("create event failed",
			logger.F("request_id", requestID),
			logger.F("venue", body.VenueID),
			logger.F("promoter", body.CreatedBy),
			logger.F("kind", logger.KindSystem),
			logger.F("error", err),
		)
		return e.InternalServerError("create event failed", err)
	}

	h.log.Info("event created",
		logger.F("request_id", requestID),
		logger.F("event", ev.ID),
		logger.F("name", ev.Name),
		logger.F("venue", ev.VenueID),
	)
	return e.JSON(201, ev)
}

func (h *Hooks) handleListEventsByPromoter(e *core.RequestEvent) error {
	requestID := e.Request.Header.Get("X-Request-Id")
	promoterID := promoter.ID(e.Request.PathValue("promoterId"))

	events, err := h.listEventByPromoter.Execute(e.Request.Context(), promoterID)
	if err != nil {
		h.log.Warn("list events by promoter failed",
			logger.F("request_id", requestID),
			logger.F("promoter", promoterID),
			logger.F("kind", logger.KindNotFound),
			logger.F("error", err),
		)
		return e.NotFoundError("promoter not found", err)
	}

	h.log.Info("listed events by promoter",
		logger.F("request_id", requestID),
		logger.F("promoter", promoterID),
		logger.F("count", len(events)),
	)
	return e.JSON(200, events)
}
