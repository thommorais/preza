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

		col := core.NewBaseCollection("promoters")

		col.Fields.Add(&core.RelationField{
			Name:          "user_id",
			CollectionId:  usersCol.Id,
			Required:      true,
			MaxSelect:     1,
			CascadeDelete: true,
		})
		col.Fields.Add(&core.TextField{
			Name:     "name",
			Required: true,
		})
		col.Fields.Add(&core.TextField{
			Name: "bio",
		})
		col.Fields.Add(&core.TextField{
			Name: "instagram",
		})
		col.Fields.Add(&core.SelectField{
			Name:      "type",
			Required:  true,
			Values:    []string{"venue", "independent"},
			MaxSelect: 1,
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
		col, err := app.FindCollectionByNameOrId("promoters")
		if err != nil {
			return err
		}
		return app.Delete(col)
	})
}
