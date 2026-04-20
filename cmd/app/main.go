package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"tournament_manager/config"
	"tournament_manager/internal/db"
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
	playerHandler := player.NewHandler(playerService, tournamentService)
	tournamentHandler.SetPlayerListFunc(func(tournamentID string) (any, error) {
		return playerService.List(tournamentID)
	})

	mux := http.NewServeMux()

	// static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// routes
	mux.HandleFunc("/", tournamentHandler.ListHandler)
	mux.HandleFunc("/tournaments/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/players") {
			playerHandler.ByTournamentHandler(w, r)
			return
		}

		tournamentHandler.ByIDHandler(w, r)
	})
	mux.HandleFunc("/tournaments", tournamentHandler.CreateHandler)

	// enrich logs by using gorilla handlers
	loggedMux := handlers.LoggingHandler(os.Stdout, mux)
	methodWrapped := middleware.MethodOverride(loggedMux)
	i18nWrapped := i18nhelper.I18nMiddleware(bundle, "en")(methodWrapped)

	if err := http.ListenAndServe(":8080", i18nWrapped); err != nil {
		log.Fatal(err)
	}
}
