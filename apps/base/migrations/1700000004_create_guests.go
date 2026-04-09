package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		usersCol, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		col := core.NewBaseCollection("guests")

		col.Fields.Add(&core.TextField{
			Name:     "name",
			Required: true,
		})
		col.Fields.Add(&core.TextField{
			Name: "phone",
		})
		col.Fields.Add(&core.TextField{
			Name: "email",
		})
		col.Fields.Add(&core.TextField{
			Name: "notes",
		})
		col.Fields.Add(&core.RelationField{
			Name:          "created_by",
			CollectionId:  usersCol.Id,
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

		return app.Save(col)
	}, func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("guests")
		if err != nil {
			return err
		}
		return app.Delete(col)
	})
}
