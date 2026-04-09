package migrations

import (
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		usersCol, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		promotersCol, err := app.FindCollectionByNameOrId("promoters")
		if err != nil {
			return err
		}
		venuesCol, err := app.FindCollectionByNameOrId("venues")
		if err != nil {
			return err
		}
		promoterVenuesCol, err := app.FindCollectionByNameOrId("promoter_venues")
		if err != nil {
			return err
		}
		guestsCol, err := app.FindCollectionByNameOrId("guests")
		if err != nil {
			return err
		}
		eventsCol, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}

		promoterNames := []string{
			"Lucas Ferreira", "Gabriel Santos", "Matheus Oliveira", "Pedro Alves",
			"Rafael Costa", "Bruno Souza", "Gustavo Lima", "Felipe Martins",
			"Thiago Rocha", "Eduardo Carvalho", "Diego Mendes", "Henrique Nunes",
			"Caio Ribeiro", "Vinicius Gomes", "Andre Pereira", "Rodrigo Barbosa",
			"Leonardo Silva", "Marcelo Araujo", "Renato Cunha", "Igor Cardoso",
		}
		promoterTypes := []string{"independent", "venue"}
		bioTemplates := []string{
			"Promotor há %d anos no cenário noturno paulistano.",
			"Especialista em eventos exclusivos e experiências premium.",
			"Conectando pessoas às melhores festas do Brasil.",
			"Criador de experiências únicas no mundo da noite.",
			"Referência em eventos de alto nível na cena carioca.",
		}
		instagramHandles := []string{
			"@promotor.lk", "@noite.gb", "@festas.mo", "@vip.pa",
			"@events.rc", "@clube.bs", "@night.gl", "@party.fm",
			"@tr.thiago", "@ec.edu", "@dm.diego", "@hn.henrique",
			"@cr.caio", "@vg.vini", "@ap.andre", "@rb.rodrigo",
			"@ls.leo", "@ma.marcelo", "@rc.renato", "@ic.igor",
		}

		userIDs := make([]string, 0, len(promoterNames))
		promoterIDs := make([]string, 0, len(promoterNames))

		for i, name := range promoterNames {
			email := fmt.Sprintf("promotor%d@preza.dev", i+1)

			userRec := core.NewRecord(usersCol)
			userRec.Set("email", email)
			userRec.Set("name", name)
			userRec.SetPassword("Preza@12345")
			userRec.Set("verified", true)
			if err := app.Save(userRec); err != nil {
				return fmt.Errorf("save user %s: %w", name, err)
			}
			userIDs = append(userIDs, userRec.Id)

			bio := fmt.Sprintf(bioTemplates[i%len(bioTemplates)], (i%8)+3)
			promoRec := core.NewRecord(promotersCol)
			promoRec.Set("user_id", userRec.Id)
			promoRec.Set("name", name)
			promoRec.Set("bio", bio)
			promoRec.Set("instagram", instagramHandles[i%len(instagramHandles)])
			promoRec.Set("type", promoterTypes[i%len(promoterTypes)])
			if err := app.Save(promoRec); err != nil {
				return fmt.Errorf("save promoter %s: %w", name, err)
			}

			promoterIDs = append(promoterIDs, promoRec.Id)
		}

		venueNames := []string{
			"Club Noir", "Terrace 54", "Paradiso SP", "The Vault",
			"Neon Garden", "Loft Ipanema", "Basement BH", "Arena Prime",
			"Casa Tulum", "Sky Lounge",
		}
		venueCities := []string{
			"São Paulo", "São Paulo", "São Paulo", "Rio de Janeiro",
			"São Paulo", "Rio de Janeiro", "Belo Horizonte", "São Paulo",
			"São Paulo", "Rio de Janeiro",
		}
		venueAddresses := []string{
			"Rua Augusta, 2345", "Av. Paulista, 1000", "Rua Oscar Freire, 800",
			"Av. Atlântica, 3500", "Rua Haddock Lobo, 1200", "Rua Visconde de Pirajá, 414",
			"Av. do Contorno, 6000", "Rua Pamplona, 900", "Alameda Lorena, 1500",
			"Av. Vieira Souto, 220",
		}
		capacities := []int{500, 300, 800, 1200, 400, 250, 600, 1500, 350, 700}
		inviteLimits := []int{50, 30, 80, 120, 40, 25, 60, 150, 35, 70}

		venueIDs := make([]string, 0, len(venueNames))

		for i, name := range venueNames {
			venueRec := core.NewRecord(venuesCol)
			venueRec.Set("promoter_id", promoterIDs[i%len(promoterIDs)])
			venueRec.Set("name", name)
			venueRec.Set("city", venueCities[i])
			venueRec.Set("address", venueAddresses[i])
			venueRec.Set("capacity", capacities[i])
			if err := app.Save(venueRec); err != nil {
				return fmt.Errorf("save venue %s: %w", name, err)
			}
			venueIDs = append(venueIDs, venueRec.Id)
		}

		// Associate promoters to venues - each venue gets 2-3 promoters
		for i, venueID := range venueIDs {
			for j := range 3 {
				promoterIdx := (i*3 + j) % len(promoterIDs)
				pvRec := core.NewRecord(promoterVenuesCol)
				pvRec.Set("venue_id", venueID)
				pvRec.Set("promoter_id", promoterIDs[promoterIdx])
				pvRec.Set("invite_limit", inviteLimits[i])
				if err := app.Save(pvRec); err != nil {
					// unique constraint violation - skip duplicate
					continue
				}
			}
		}

		guestFirstNames := []string{
			"Ana", "Beatriz", "Camila", "Daniela", "Eduarda",
			"Fernanda", "Gabriela", "Helena", "Isabela", "Juliana",
			"Karen", "Larissa", "Marina", "Natalia", "Olivia",
			"Patricia", "Rafaela", "Sabrina", "Tatiane", "Vanessa",
			"Carlos", "Daniel", "Eduardo", "Felipe", "Guilherme",
			"Henrique", "Igor", "Joao", "Leonardo", "Marcos",
			"Nicolas", "Otavio", "Paulo", "Rodrigo", "Samuel",
			"Thiago", "Ulisses", "Victor", "William", "Yuri",
		}
		guestLastNames := []string{
			"Almeida", "Barbosa", "Carvalho", "Costa", "Dias",
			"Ferreira", "Gomes", "Lopes", "Martins", "Nascimento",
			"Oliveira", "Pereira", "Ribeiro", "Santos", "Silva",
			"Souza", "Teixeira", "Vieira", "Xavier", "Zanetti",
		}
		domainSuffixes := []string{
			"gmail.com", "hotmail.com", "yahoo.com.br", "outlook.com", "uol.com.br",
		}
		phonePrefixes := []string{
			"11", "21", "31", "41", "51", "61", "71", "81", "85", "92",
		}
		noteTemplates := []string{
			"VIP frequente.", "Amigo do promotor.", "Cliente especial.",
			"Convidado da casa.", "Lista A.", "", "", "", "", "",
		}

		for i := range 100 {
			firstName := guestFirstNames[i%len(guestFirstNames)]
			lastName := guestLastNames[(i/len(guestFirstNames))%len(guestLastNames)]
			fullName := fmt.Sprintf("%s %s", firstName, lastName)
			email := fmt.Sprintf("%s.%s%d@%s",
				sanitizeName(firstName),
				sanitizeName(lastName),
				i+1,
				domainSuffixes[i%len(domainSuffixes)],
			)
			phone := fmt.Sprintf("+55%s9%04d%04d",
				phonePrefixes[i%len(phonePrefixes)],
				(i*7+1337)%10000,
				(i*13+4242)%10000,
			)

			guestRec := core.NewRecord(guestsCol)
			guestRec.Set("name", fullName)
			guestRec.Set("email", email)
			guestRec.Set("phone", phone)
			guestRec.Set("notes", noteTemplates[i%len(noteTemplates)])
			guestRec.Set("created_by", userIDs[i%len(userIDs)])

			if err := app.Save(guestRec); err != nil {
				return fmt.Errorf("save guest %s: %w", fullName, err)
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
		usersCol, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		for i := range 20 {
			email := fmt.Sprintf("promotor%d@preza.dev", i+1)
			rec, err := app.FindAuthRecordByEmail(usersCol.Id, email)
			if err != nil {
				continue
			}
			_ = app.Delete(rec)
		}
		return nil
	})
}

func sanitizeName(s string) string {
	result := make([]byte, 0, len(s))
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			result = append(result, byte(c+32))
		} else if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
