package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		eventsCol, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}

		promotersCol, err := app.FindCollectionByNameOrId("promoters")
		if err != nil {
			return err
		}

		guestsCol, err := app.FindCollectionByNameOrId("guests")
		if err != nil {
			return err
		}

		col := core.NewBaseCollection("invitations")

		col.Fields.Add(&core.RelationField{
			Name:          "event_id",
			CollectionId:  eventsCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
		})
		col.Fields.Add(&core.RelationField{
			Name:          "created_by",
			CollectionId:  promotersCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
		})
		col.Fields.Add(&core.RelationField{
			Name:          "guest_id",
			CollectionId:  guestsCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
		})
		col.Fields.Add(&core.SelectField{
			Name:      "status",
			Values:    []string{"pending", "confirmed", "declined", "used"},
			MaxSelect: 1,
		})
		col.Fields.Add(&core.TextField{
			Name:     "token",
			Required: true,
		})
		col.Fields.Add(&core.NumberField{
			Name: "plus_ones",
		})
		col.Fields.Add(&core.JSONField{
			Name: "plus_one_names",
		})
		col.Fields.Add(&core.DateField{
			Name: "rsvp_at",
		})
		col.Fields.Add(&core.DateField{
			Name: "checked_in_at",
		})
		col.Fields.Add(&core.TextField{
			Name: "notes",
		})
		col.Fields.Add(&core.AutodateField{
			Name:     "created",
			OnCreate: true,
		})
		col.Fields.Add(&core.AutodateField{
			Name:     "updated",
			OnCreate: true,
			OnUpdate: true,
		})

		col.AddIndex("idx_invitations_token_unique", true, "token", "")

		return app.Save(col)
	}, func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("invitations")
		if err != nil {
			return err
		}
		return app.Delete(col)
	})
}
