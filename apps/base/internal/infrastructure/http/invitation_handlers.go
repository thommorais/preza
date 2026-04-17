package http

import (
	"encoding/json"
	"errors"
	"preza/internal/application/invitation"
	domainInvitation "preza/internal/domain/invitation"
	"preza/internal/domain/event"
	"preza/internal/domain/guest"
	"preza/internal/domain/logger"
	"preza/internal/domain/promoter"

	"github.com/pocketbase/pocketbase/core"
)

func (h *Hooks) handleCheckIn(e *core.RequestEvent) error {
	token := e.Request.PathValue("token")
	requestID := e.Request.Header.Get("X-Request-Id")

	inv, err := h.checkIn.Execute(e.Request.Context(), token)
	if err != nil {
		if errors.Is(err, domainInvitation.ErrNotFound) || errors.Is(err, domainInvitation.ErrInvalidToken) {
			h.log.Warn("check-in failed: token not found",
				logger.F("request_id", requestID),
				logger.F("token", token),
				logger.F("kind", logger.KindNotFound),
				logger.F("error", err),
			)
			return e.NotFoundError("invitation not found", err)
		}
		if errors.Is(err, domainInvitation.ErrAlreadyCheckedIn) {
			h.log.Warn("check-in failed: already checked in",
				logger.F("request_id", requestID),
				logger.F("token", token),
				logger.F("kind", logger.KindValidation),
				logger.F("error", err),
			)
			return e.JSON(409, map[string]any{"code": 409, "message": err.Error()})
		}
		h.log.Error("check-in failed",
			logger.F("request_id", requestID),
			logger.F("token", token),
			logger.F("kind", logger.KindSystem),
			logger.F("error", err),
		)
		return e.InternalServerError("check-in failed", err)
	}

	h.log.Info("guest checked in",
		logger.F("request_id", requestID),
		logger.F("invitation", inv.ID),
		logger.F("event", inv.EventID),
	)
	return e.JSON(200, inv)
}

func (h *Hooks) handleCreateInvitation(e *core.RequestEvent) error {
	requestID := e.Request.Header.Get("X-Request-Id")

	var body struct {
		EventID   string `json:"event_id"`
		CreatedBy string `json:"created_by"`
		GuestID   string `json:"guest_id"`
		PlusOnes  int    `json:"plus_ones"`
		Notes     string `json:"notes"`
	}

	if err := json.NewDecoder(e.Request.Body).Decode(&body); err != nil {
		h.log.Warn("invalid create invitation body",
			logger.F("request_id", requestID),
			logger.F("kind", logger.KindValidation),
			logger.F("error", err),
		)
		return e.BadRequestError("invalid request body", err)
	}

	inv, err := h.createInvitation.Execute(e.Request.Context(), invitation.CreateInvitationInput{
		EventID:   event.ID(body.EventID),
		CreatedBy: promoter.ID(body.CreatedBy),
		GuestID:   guest.ID(body.GuestID),
		PlusOnes:  body.PlusOnes,
		Notes:     body.Notes,
	})
	if err != nil {
		if errors.Is(err, invitation.ErrEventNotActive) {
			h.log.Warn("create invitation failed: event not active",
				logger.F("request_id", requestID),
				logger.F("event", body.EventID),
				logger.F("kind", logger.KindValidation),
				logger.F("error", err),
			)
			return e.JSON(422, map[string]any{"code": 422, "message": err.Error()})
		}
		h.log.Error("create invitation failed",
			logger.F("request_id", requestID),
			logger.F("event", body.EventID),
			logger.F("guest", body.GuestID),
			logger.F("kind", logger.KindSystem),
			logger.F("error", err),
		)
		return e.InternalServerError("create invitation failed", err)
	}

	h.log.Info("invitation created",
		logger.F("request_id", requestID),
		logger.F("invitation", inv.ID),
		logger.F("event", inv.EventID),
		logger.F("guest", inv.GuestID),
	)
	return e.JSON(201, inv)
}

func (h *Hooks) handleListInvitationsByEvent(e *core.RequestEvent) error {
	requestID := e.Request.Header.Get("X-Request-Id")
	eventID := event.ID(e.Request.PathValue("eventId"))

	invitations, err := h.listInvByEvent.Execute(e.Request.Context(), eventID)
	if err != nil {
		h.log.Warn("list invitations by event failed",
			logger.F("request_id", requestID),
			logger.F("event", eventID),
			logger.F("kind", logger.KindNotFound),
			logger.F("error", err),
		)
		return e.NotFoundError("event not found", err)
	}

	h.log.Info("listed invitations by event",
		logger.F("request_id", requestID),
		logger.F("event", eventID),
		logger.F("count", len(invitations)),
	)
	return e.JSON(200, invitations)
}

func (h *Hooks) handleListInvitationsByPromoter(e *core.RequestEvent) error {
	requestID := e.Request.Header.Get("X-Request-Id")
	promoterID := promoter.ID(e.Request.PathValue("promoterId"))

	invitations, err := h.listInvByPromoter.Execute(e.Request.Context(), promoterID)
	if err != nil {
		h.log.Warn("list invitations by promoter failed",
			logger.F("request_id", requestID),
			logger.F("promoter", promoterID),
			logger.F("kind", logger.KindNotFound),
			logger.F("error", err),
		)
		return e.NotFoundError("promoter not found", err)
	}

	h.log.Info("listed invitations by promoter",
		logger.F("request_id", requestID),
		logger.F("promoter", promoterID),
		logger.F("count", len(invitations)),
	)
	return e.JSON(200, invitations)
}
