package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		venuesCol, err := app.FindCollectionByNameOrId("venues")
		if err != nil {
			return err
		}

		promotersCol, err := app.FindCollectionByNameOrId("promoters")
		if err != nil {
			return err
		}

		col := core.NewBaseCollection("events")

		col.Fields.Add(&core.RelationField{
			Name:          "venue_id",
			CollectionId:  venuesCol.Id,
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
		col.Fields.Add(&core.TextField{
			Name:     "name",
			Required: true,
		})
		col.Fields.Add(&core.DateField{
			Name:     "date",
			Required: true,
		})
		col.Fields.Add(&core.NumberField{
			Name: "max_capacity",
		})
		col.Fields.Add(&core.SelectField{
			Name:      "status",
			Values:    []string{"draft", "active", "closed"},
			MaxSelect: 1,
		})
		col.Fields.Add(&core.FileField{
			Name:      "cover_image",
			MaxSelect: 1,
			MimeTypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
		})
		col.Fields.Add(&core.TextField{
			Name: "description",
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
		col, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}
		return app.Delete(col)
	})
}
