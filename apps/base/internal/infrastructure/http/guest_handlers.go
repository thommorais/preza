package http

import (
	"preza/internal/domain/logger"
	"preza/internal/domain/promoter"

	"github.com/pocketbase/pocketbase/core"
)

func (h *Hooks) handleListGuestsByPromoter(e *core.RequestEvent) error {
	requestID := e.Request.Header.Get("X-Request-Id")
	promoterID := promoter.ID(e.Request.PathValue("promoterId"))

	guests, err := h.listGuestByPromoter.Execute(e.Request.Context(), promoterID)
	if err != nil {
		h.log.Warn("list guests by promoter failed",
			logger.F("request_id", requestID),
			logger.F("promoter", promoterID),
			logger.F("kind", logger.KindNotFound),
			logger.F("error", err),
		)
		return e.NotFoundError("promoter not found", err)
	}

	h.log.Info("listed guests by promoter",
		logger.F("request_id", requestID),
		logger.F("promoter", promoterID),
		logger.F("count", len(guests)),
	)
	return e.JSON(200, guests)
}
