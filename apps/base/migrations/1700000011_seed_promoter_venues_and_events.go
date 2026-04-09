package migrations

import (
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		promoterVenuesCol, err := app.FindCollectionByNameOrId("promoter_venues")
		if err != nil {
			return err
		}
		eventsCol, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}

		// Load existing promoters and venues
		promoterRecords, err := app.FindRecordsByFilter("promoters", "id != ''", "-created", 20, 0)
		if err != nil {
			return fmt.Errorf("load promoters: %w", err)
		}
		venueRecords, err := app.FindRecordsByFilter("venues", "id != ''", "-created", 10, 0)
		if err != nil {
			return fmt.Errorf("load venues: %w", err)
		}

		promoterIDs := make([]string, 0, len(promoterRecords))
		for _, r := range promoterRecords {
			promoterIDs = append(promoterIDs, r.Id)
		}
		venueIDs := make([]string, 0, len(venueRecords))
		capacities := make([]int, 0, len(venueRecords))
		for _, r := range venueRecords {
			venueIDs = append(venueIDs, r.Id)
			capacities = append(capacities, int(r.GetFloat("capacity")))
		}

		inviteLimits := []int{50, 30, 80, 120, 40, 25, 60, 150, 35, 70}

		// Each venue gets 3 promoters linked
		for i, venueID := range venueIDs {
			for j := range 3 {
				promoterIdx := (i*3 + j) % len(promoterIDs)
				pvRec := core.NewRecord(promoterVenuesCol)
				pvRec.Set("venue_id", venueID)
				pvRec.Set("promoter_id", promoterIDs[promoterIdx])
				pvRec.Set("invite_limit", inviteLimits[i%len(inviteLimits)])
				if err := app.Save(pvRec); err != nil {
					// skip unique constraint violations
					continue
				}
			}
		}

		now := time.Now()

		// 5 future events
		futureEventNames := []string{
			"Réveillon Preza 2026", "Carnaval VIP Experience",
			"Festival Noite Paulistana", "Club Noir Opening Night", "Terrace Sessions Vol. 3",
		}
		futureEventDescs := []string{
			"A maior virada do ano com open bar premium e DJs internacionais.",
			"Carnaval exclusivo para convidados selecionados, blocos privativos e muito samba.",
			"Três palcos, vinte atrações e a melhor noite do ano em São Paulo.",
			"Abertura oficial da nova temporada do Club Noir com line-up surpresa.",
			"Tarde e noite na terraça com sunset sessions e DJs residentes.",
		}
		futureDaysAhead := []int{7, 21, 45, 60, 90}

		for i := range 5 {
			eventDate := now.AddDate(0, 0, futureDaysAhead[i])
			eventRec := core.NewRecord(eventsCol)
			eventRec.Set("venue_id", venueIDs[i%len(venueIDs)])
			eventRec.Set("created_by", promoterIDs[i%len(promoterIDs)])
			eventRec.Set("name", futureEventNames[i])
			eventRec.Set("description", futureEventDescs[i])
			eventRec.Set("date", eventDate.Format("2006-01-02 15:04:05"))
			eventRec.Set("max_capacity", capacities[i%len(capacities)])
			eventRec.Set("status", "active")
			if err := app.Save(eventRec); err != nil {
				return fmt.Errorf("save future event %s: %w", futureEventNames[i], err)
			}
		}

		// 20 past events
		pastEventNames := []string{
			"Ano Novo Preza 2025", "Carnaval Block Party", "Dia dos Namorados VIP",
			"Festa Junina Premium", "Aniversário Club Noir", "Halloween Noir",
			"Black Friday Night", "Natal Preza", "Reveillon 2024",
			"São Paulo Fashion Night", "Copa do Mundo Special", "Dia da Música",
			"Festival de Inverno", "Semana do Promotor", "Open Bar Sábado",
			"Sunset Terrace", "Rave da Madrugada", "Domingo de Funk",
			"Pagode VIP", "Sertanejo Premium",
		}
		pastEventDescs := []string{
			"Virada do ano com open bar e DJs residentes.",
			"Bloco privativo de carnaval no coração de SP.",
			"Jantar harmonizado e festa exclusiva para casais.",
			"Forró, quadrilha e comida típica em clima de São João.",
			"Cinco anos de história com line-up especial.",
			"A festa de Halloween mais assustadora de SP.",
			"Black Friday virou black night - festa com os melhores preços.",
			"Ceia de Natal com open food e open bar.",
			"Contagem regressiva do ano com fogos e DJ.",
			"Noite de moda e música eletrônica no coração da cidade.",
			"Final da Copa transmitida ao vivo com festa temática.",
			"Celebração do Dia Internacional da Música.",
			"Festival de jazz e blues na temporada de inverno.",
			"Semana especial com festas todos os dias.",
			"Open bar premium com DJs nacionais.",
			"Sunset com vista para a cidade e DJ ao vivo.",
			"Rave de madrugada com line-up underground.",
			"Baile funk na laje com os melhores MCs.",
			"Pagode VIP com grandes artistas.",
			"Sertanejo universitário com duplas do momento.",
		}

		for i := range 20 {
			daysAgo := (i+1)*15 + (i * 7)
			eventDate := now.AddDate(0, 0, -daysAgo)
			eventRec := core.NewRecord(eventsCol)
			eventRec.Set("venue_id", venueIDs[i%len(venueIDs)])
			eventRec.Set("created_by", promoterIDs[i%len(promoterIDs)])
			eventRec.Set("name", pastEventNames[i])
			eventRec.Set("description", pastEventDescs[i])
			eventRec.Set("date", eventDate.Format("2006-01-02 15:04:05"))
			eventRec.Set("max_capacity", capacities[i%len(capacities)])
			eventRec.Set("status", "closed")
			if err := app.Save(eventRec); err != nil {
				return fmt.Errorf("save past event %s: %w", pastEventNames[i], err)
			}
		}

		return nil
	}, func(app core.App) error {
		eventsCol, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}
		records, err := app.FindRecordsByFilter("events", "id != ''", "-created", 25, 0)
		if err != nil {
			return err
		}
		for _, r := range records {
			_ = app.Delete(r)
		}

		pvRecords, err := app.FindRecordsByFilter("promoter_venues", "id != ''", "-created", 30, 0)
		if err != nil {
			return err
		}
		for _, r := range pvRecords {
			_ = app.Delete(r)
		}

		_ = eventsCol
		return nil
	})
}
