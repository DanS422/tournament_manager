package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"tournament_manager/config"
	"tournament_manager/internal/attendant"
	"tournament_manager/internal/db"
	"tournament_manager/internal/discipline"
	i18nhelper "tournament_manager/internal/i18n"
	"tournament_manager/internal/middleware"
	"tournament_manager/internal/player"
	"tournament_manager/internal/tournament"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/handlers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	cfg := config.Load()
	dsn, found := strings.CutPrefix(cfg.DatabaseURL, "sqlite://")

	if !found {
		log.Fatal("Database URL must start with sqlite://")
	}

	dbConn, err := db.Open(dsn)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to database %s", cfg.DatabaseURL)

	switch os.Getenv("MIGRATE") {
	case "up":
		if err := db.MigrateUp(dbConn); err != nil {
			log.Fatal(err)
		}
		return
	case "down":
		if err := db.MigrateDown(dbConn); err != nil {
			log.Fatal(err)
		}
		return
	case "reset":
		if err := db.MigrateReset(dbConn); err != nil {
			log.Fatal(err)
		}
		return
	}

	// normal app startup

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("locales/en.toml")
	bundle.LoadMessageFile("locales/de.toml")

	tournamentRepo := tournament.NewRepository(dbConn)
	tournamentService := tournament.NewService(tournamentRepo)
	tournamentHandler := tournament.NewHandler(tournamentService)

	playerRepo := player.NewRepository(dbConn)
	playerService := player.NewService(playerRepo)
	playerHandler := player.NewHandler(playerService)

	attendantRepo := attendant.NewRepository(dbConn)
	attendantService := attendant.NewService(attendantRepo)
	attendantHandler := attendant.NewHandler(attendantService)

	disciplineRepo := discipline.NewRepository(dbConn)
	disciplineService := discipline.NewService(disciplineRepo)
	disciplineHandler := discipline.NewHandler(disciplineService)

	mux := http.NewServeMux()

	// static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// routes

	// Tournaments
	mux.HandleFunc("/", tournamentHandler.ListHandler)
	mux.HandleFunc("/tournaments", tournamentHandler.CreateHandler)
	mux.HandleFunc("/tournaments/", func(w http.ResponseWriter, r *http.Request) {
		tournamentHandler.ByIDHandler(w, r)
	})

	// Attendants
	mux.HandleFunc("GET /tournaments/{tournament_id}/attendants", func(w http.ResponseWriter, r *http.Request) {
		attendantHandler.ListHandler(w, r, playerService)
	})
	mux.HandleFunc("POST /tournaments/{tournament_id}/attendants", attendantHandler.CreateHandler)
	mux.HandleFunc("DELETE /tournaments/{tournament_id}/attendants/{attendant_id}", attendantHandler.DeleteHandler)
	mux.HandleFunc("GET /tournaments/{tournament_id}/attendants/{attendant_id}", attendantHandler.ShowHandler)

	// Discipline
	mux.HandleFunc("GET /tournaments/{tournament_id}/disciplines", disciplineHandler.ListHandler)
	mux.HandleFunc("POST /tournaments/{tournament_id}/disciplines", disciplineHandler.CreateHandler)
	mux.HandleFunc("GET /tournaments/{tournament_id}/disciplines/{discipline_id}", disciplineHandler.ShowHandler)
	mux.HandleFunc("DELETE /tournaments/{tournament_id}/disciplines/{discipline_id}", disciplineHandler.DeleteHandler)
	mux.HandleFunc("PATCH /tournaments/{tournament_id}/disciplines/{discipline_id}", disciplineHandler.UpdateHandler)

	// Players
	mux.HandleFunc("/players", playerHandler.PlayersHandler)
	mux.HandleFunc("/players/", playerHandler.PlayersHandler)

	// enrich logs by using gorilla handlers
	loggedMux := handlers.LoggingHandler(os.Stdout, mux)
	methodWrapped := middleware.MethodOverride(loggedMux)
	i18nWrapped := i18nhelper.I18nMiddleware(bundle, "en")(methodWrapped)

	if err := http.ListenAndServe(":8080", i18nWrapped); err != nil {
		log.Fatal(err)
	}
}
