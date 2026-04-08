package main

import (
	"net/http"
	"tournament_manager/internal/middleware"
	"tournament_manager/internal/tournament"
)

func main() {
	repo := tournament.NewRepository()
	service := tournament.NewService(repo)
	handler := tournament.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.ListHandler)
	mux.HandleFunc("/tournaments/", handler.ByIDHandler)
	mux.HandleFunc("/tournaments", handler.CreateHandler)

	wrapped := middleware.MethodOverride(mux)

	http.ListenAndServe(":8080", wrapped)
}
