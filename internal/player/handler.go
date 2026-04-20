package player

import (
	"html/template"
	"net/http"
	"strings"
	tmpl "tournament_manager/internal/tmpl"
	"tournament_manager/internal/tournament"
	"tournament_manager/internal/validation"

	"github.com/google/uuid"
)

type TournamentService interface {
	Show(id string) (tournament.Tournament, error)
}

type Handler struct {
	service           ServiceInterface
	tournamentService TournamentService
	templates         map[string]*template.Template
}

func NewHandler(s ServiceInterface, ts TournamentService) *Handler {
	t := make(map[string]*template.Template)
	t["tournament"] = template.Must(template.New("").
		Funcs(tmpl.FuncMap()).
		ParseFiles(
			"templates/base.html",
			"templates/tournament.html",
		))

	return &Handler{
		service:           s,
		tournamentService: ts,
		templates:         t,
	}
}

func (h *Handler) ByTournamentHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID, playerID, ok := parsePlayerPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if _, err := uuid.Parse(tournamentID); err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusNotFound)
		return
	}

	if playerID != "" {
		if _, err := uuid.Parse(playerID); err != nil {
			http.Error(w, "Invalid player ID", http.StatusNotFound)
			return
		}
	}

	switch r.Method {
	case http.MethodPost:
		if playerID != "" {
			http.NotFound(w, r)
			return
		}
		h.createHandler(w, r, tournamentID)
	case http.MethodPatch, http.MethodPut:
		if playerID == "" {
			http.NotFound(w, r)
			return
		}
		h.updateHandler(w, r, tournamentID, playerID)
	case http.MethodDelete:
		if playerID == "" {
			http.NotFound(w, r)
			return
		}
		h.deleteHandler(w, r, tournamentID, playerID)
	default:
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createHandler(w http.ResponseWriter, r *http.Request, tournamentID string) {
	p, errs := playerFromRequest(r, tournamentID, "")
	if len(errs) > 0 {
		h.renderTournament(w, r, tournamentID, p, errs)
		return
	}

	if _, err := h.service.Create(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID, http.StatusSeeOther)
}

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request, tournamentID, playerID string) {
	p, errs := playerFromRequest(r, tournamentID, playerID)
	if len(errs) > 0 {
		h.renderTournament(w, r, tournamentID, p, errs)
		return
	}

	if err := h.service.Update(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID, http.StatusSeeOther)
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request, tournamentID, playerID string) {
	if err := h.service.Delete(playerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID, http.StatusSeeOther)
}

func (h *Handler) renderTournament(w http.ResponseWriter, r *http.Request, tournamentID string, p Player, errs map[string]string) {
	tour, err := h.tournamentService.Show(tournamentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	players, err := h.service.List(tournamentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["tournament"], map[string]interface{}{
		"Errors":     errs,
		"Player":     p,
		"Players":    players,
		"Tournament": tour,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func playerFromRequest(r *http.Request, tournamentID, playerID string) (Player, map[string]string) {
	p := Player{
		ID:           playerID,
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
		Gender:       r.FormValue("gender"),
		TournamentID: tournamentID,
	}

	errs := map[string]string{}
	for field, msg := range validation.ValidateStruct(p) {
		errs[field] = msg
	}

	return p, errs
}

func parsePlayerPath(path string) (string, string, bool) {
	rest := strings.TrimPrefix(path, "/tournaments/")
	parts := strings.Split(rest, "/")

	if len(parts) == 2 && parts[1] == "players" {
		return parts[0], "", parts[0] != ""
	}

	if len(parts) == 3 && parts[1] == "players" {
		return parts[0], parts[2], parts[0] != "" && parts[2] != ""
	}

	return "", "", false
}
