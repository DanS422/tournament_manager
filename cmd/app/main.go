package main

import (
	"net/http"
	"os"
	"tournament_manager/internal/db"
	i18nhelper "tournament_manager/internal/i18n"
	"tournament_manager/internal/middleware"
	"tournament_manager/internal/tournament"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/handlers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	switch os.Getenv("MIGRATE") {
	case "up":
		db.MigrateUp()
		return
	case "down":
		db.MigrateDown()
		return
	case "reset":
		db.MigrateReset()
		return
	}

	// normal app startup

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("locales/en.toml")
	bundle.LoadMessageFile("locales/de.toml")

	repo, err := tournament.NewRepository()

	if err != nil {
		panic("database connection failed")
	}
	service := tournament.NewService(repo)
	handler := tournament.NewHandler(service)

	mux := http.NewServeMux()

	// static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// routes
	mux.HandleFunc("/", handler.ListHandler)
	mux.HandleFunc("/tournaments/", handler.ByIDHandler)
	mux.HandleFunc("/tournaments", handler.CreateHandler)

	// enrich logs by using gorilla handlers
	loggedMux := handlers.LoggingHandler(os.Stdout, mux)
	methodWrapped := middleware.MethodOverride(loggedMux)
	i18nWrapped := i18nhelper.I18nMiddleware(bundle, "en")(methodWrapped)

	http.ListenAndServe(":8080", i18nWrapped)
}
