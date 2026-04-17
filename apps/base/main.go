package main

import (
	_ "embed"
	"log"
	"os"
	"strings"

	appevent "preza/internal/application/event"
	appguest "preza/internal/application/guest"
	appinvitation "preza/internal/application/invitation"
	httpadapter "preza/internal/infrastructure/http"
	inflogger "preza/internal/infrastructure/logger"
	pb "preza/internal/infrastructure/pocketbase"
	_ "preza/migrations"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

//go:embed openapi.yaml
var openAPISpec []byte

func main() {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	logger := inflogger.New()

	invitationRepo := pb.NewInvitationRepository(app, logger)
	eventRepo := pb.NewEventRepository(app, logger)
	venueRepo := pb.NewVenueRepository(app, logger)
	promoterRepo := pb.NewPromoterRepository(app, logger)
	guestRepo := pb.NewGuestRepository(app, logger)

	hooks := httpadapter.NewHooks(httpadapter.HooksConfig{
		Log:                 logger,
		OpenAPISpec:         openAPISpec,
		CheckIn:             appinvitation.NewCheckInGuestUseCase(invitationRepo),
		CreateInvitation:    appinvitation.NewCreateInvitationUseCase(invitationRepo, eventRepo),
		ListInvByEvent:      appinvitation.NewListByEventUseCase(invitationRepo, eventRepo),
		ListInvByPromoter:   appinvitation.NewListByPromoterUseCase(invitationRepo, promoterRepo),
		CreateEvent:         appevent.NewCreateEventUseCase(eventRepo, venueRepo, promoterRepo),
		ListEventByPromoter: appevent.NewListByPromoterUseCase(eventRepo, promoterRepo),
		ListGuestByPromoter: appguest.NewListByPromoterUseCase(guestRepo, promoterRepo),
	})
	hooks.Register(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
