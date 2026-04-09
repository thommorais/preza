package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		promotersCol, err := app.FindCollectionByNameOrId("promoters")
		if err != nil {
			return err
		}

		venuesCol, err := app.FindCollectionByNameOrId("venues")
		if err != nil {
			return err
		}

		col := core.NewBaseCollection("promoter_venues")

		col.Fields.Add(&core.RelationField{
			Name:          "promoter_id",
			CollectionId:  promotersCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
		})
		col.Fields.Add(&core.RelationField{
			Name:          "venue_id",
			CollectionId:  venuesCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
		})
		col.Fields.Add(&core.NumberField{
			Name: "invite_limit",
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

		col.AddIndex("idx_promoter_venues_unique", true, "promoter_id, venue_id", "")

		return app.Save(col)
	}, func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("promoter_venues")
		if err != nil {
			return err
		}
		return app.Delete(col)
	})
}
