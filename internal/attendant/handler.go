package attendant

import (
	"html/template"
	"net/http"
	"tournament_manager/internal/player"
	"tournament_manager/internal/tmpl"
)

type Handler struct {
	service   ServiceInterface
	templates map[string]*template.Template
}

func NewHandler(s ServiceInterface) *Handler {
	t := make(map[string]*template.Template)

	pages := []string{"attendants", "attendant"}

	for _, p := range pages {
		t[p] = template.Must(template.New("").
			Funcs(tmpl.FuncMap()).
			ParseFiles(
				"templates/base.html",
				"templates/"+p+".html",
			))
	}

	return &Handler{
		service:   s,
		templates: t,
	}
}

func (h *Handler) ShowHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")
	attendantID := r.PathValue("attendant_id")

	attendant, err := h.service.Show(tournamentID, attendantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["attendant"], map[string]interface{}{
		"Attendant": attendant,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request, playService *player.Service) {
	tournamentID := r.PathValue("tournament_id")

	attendants, err := h.service.GetAll(tournamentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendantMap := make(map[string]bool)
	for _, attendant := range attendants {
		attendantMap[attendant.PlayerID] = true
	}

	players, err := playService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["attendants"], map[string]interface{}{
		"Attendants":   attendants,
		"Players":      players,
		"TournamentID": tournamentID,
		"AttendingMap": attendantMap,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")
	a := Attendant{
		TournamentID: tournamentID,
		PlayerID:     r.FormValue("player_id"),
	}

	_, err := h.service.Add(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID+"/attendants", http.StatusSeeOther)
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")
	attendantID := r.PathValue("attendant_id")

	err := h.service.Delete(tournamentID, attendantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID+"/attendants", http.StatusSeeOther)
}
