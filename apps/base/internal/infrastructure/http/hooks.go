package http

import (
	appinvitation "preza/internal/application/invitation"

	"github.com/pocketbase/pocketbase/core"
)

type Hooks struct {
	checkIn *appinvitation.CheckInGuestUseCase
	create  *appinvitation.CreateInvitationUseCase
}

func NewHooks(checkIn *appinvitation.CheckInGuestUseCase, create *appinvitation.CreateInvitationUseCase) *Hooks {
	return &Hooks{checkIn: checkIn, create: create}
}

func (h *Hooks) Register(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.POST("/api/invitations/{token}/check-in", h.handleCheckIn)
		return e.Next()
	})
}

func (h *Hooks) handleCheckIn(e *core.RequestEvent) error {
	token := e.Request.PathValue("token")

	inv, err := h.checkIn.Execute(e.Request.Context(), token)
	if err != nil {
		return e.BadRequestError("check-in failed", err)
	}

	return e.JSON(200, inv)
}
