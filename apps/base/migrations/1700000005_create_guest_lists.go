package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		guestsCol, err := app.FindCollectionByNameOrId("guests")
		if err != nil {
			return err
		}

		promotersCol, err := app.FindCollectionByNameOrId("promoters")
		if err != nil {
			return err
		}

		col := core.NewBaseCollection("guest_lists")

		col.Fields.Add(&core.RelationField{
			Name:          "guest_id",
			CollectionId:  guestsCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
		})
		col.Fields.Add(&core.RelationField{
			Name:          "promoter_id",
			CollectionId:  promotersCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
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

		col.AddIndex("idx_guest_lists_unique", true, "guest_id, promoter_id", "")

		return app.Save(col)
	}, func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("guest_lists")
		if err != nil {
			return err
		}
		return app.Delete(col)
	})
}
