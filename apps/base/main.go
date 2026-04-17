package main

import (
	"log"
	"os"
	"strings"

	appinvitation "preza/internal/application/invitation"
	httpadapter "preza/internal/infrastructure/http"
	pb "preza/internal/infrastructure/pocketbase"
	_ "preza/migrations"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	invitationRepo := pb.NewInvitationRepository(app)
	eventRepo := pb.NewEventRepository(app)
	checkIn := appinvitation.NewCheckInGuestUseCase(invitationRepo)
	createInvitation := appinvitation.NewCreateInvitationUseCase(invitationRepo, eventRepo)

	hooks := httpadapter.NewHooks(checkIn, createInvitation)
	hooks.Register(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
