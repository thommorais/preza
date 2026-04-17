package http

import (
	appevent "preza/internal/application/event"
	appguest "preza/internal/application/guest"
	appinvitation "preza/internal/application/invitation"
	"preza/internal/domain/logger"

	"github.com/pocketbase/pocketbase/core"
)

type Hooks struct {
	log                 logger.Logger
	openAPISpec         []byte
	checkIn             *appinvitation.CheckInGuestUseCase
	createInvitation    *appinvitation.CreateInvitationUseCase
	listInvByEvent      *appinvitation.ListByEventUseCase
	listInvByPromoter   *appinvitation.ListByPromoterUseCase
	createEvent         *appevent.CreateEventUseCase
	listEventByPromoter *appevent.ListByPromoterUseCase
	listGuestByPromoter *appguest.ListByPromoterUseCase
}

type HooksConfig struct {
	Log                 logger.Logger
	OpenAPISpec         []byte
	CheckIn             *appinvitation.CheckInGuestUseCase
	CreateInvitation    *appinvitation.CreateInvitationUseCase
	ListInvByEvent      *appinvitation.ListByEventUseCase
	ListInvByPromoter   *appinvitation.ListByPromoterUseCase
	CreateEvent         *appevent.CreateEventUseCase
	ListEventByPromoter *appevent.ListByPromoterUseCase
	ListGuestByPromoter *appguest.ListByPromoterUseCase
}

func NewHooks(cfg HooksConfig) *Hooks {
	return &Hooks{
		log:                 cfg.Log,
		openAPISpec:         cfg.OpenAPISpec,
		checkIn:             cfg.CheckIn,
		createInvitation:    cfg.CreateInvitation,
		listInvByEvent:      cfg.ListInvByEvent,
		listInvByPromoter:   cfg.ListInvByPromoter,
		createEvent:         cfg.CreateEvent,
		listEventByPromoter: cfg.ListEventByPromoter,
		listGuestByPromoter: cfg.ListGuestByPromoter,
	}
}

func (h *Hooks) Register(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		r := e.Router

		r.BindFunc(requestMiddleware(h.log))

		r.PATCH("/api/invitations/{token}/check-in", h.handleCheckIn)
		r.POST("/api/invitations", h.handleCreateInvitation)
		r.GET("/api/events/{eventId}/invitations", h.handleListInvitationsByEvent)
		r.GET("/api/promoters/{promoterId}/invitations", h.handleListInvitationsByPromoter)

		r.POST("/api/events", h.handleCreateEvent)
		r.GET("/api/promoters/{promoterId}/events", h.handleListEventsByPromoter)

		r.GET("/api/promoters/{promoterId}/guests", h.handleListGuestsByPromoter)

		r.GET("/api/openapi.yaml", h.handleOpenAPISpec)
		r.GET("/api/docs", h.handleSwaggerUI)

		return e.Next()
	})
}
