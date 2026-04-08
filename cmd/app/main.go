package main

import (
	"net/http"
	"os"
	"tournament_manager/internal/middleware"
	"tournament_manager/internal/tournament"

	"github.com/gorilla/handlers"
)

func main() {
	repo := tournament.NewRepository()
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
	wrapped := middleware.MethodOverride(loggedMux)

	http.ListenAndServe(":8080", wrapped)
}
