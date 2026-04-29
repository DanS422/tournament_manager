package discipline

import (
	"html/template"
	"net/http"
	"strconv"
	"tournament_manager/internal/tmpl"
	"tournament_manager/internal/validation"
)

type Handler struct {
	service   ServiceInterface
	templates map[string]*template.Template
}

func NewHandler(s ServiceInterface) *Handler {
	t := make(map[string]*template.Template)
	pages := []string{"disciplines", "discipline"}

	for _, p := range pages {
		t[p] = template.Must(template.New("").
			Funcs(tmpl.FuncMap()).
			ParseFiles(
				"templates/base.html",
				"templates/"+p+".html",
			),
		)
	}

	return &Handler{
		service:   s,
		templates: t,
	}
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")
	no_of_team_players, err := strconv.Atoi(r.FormValue("no_of_team_players"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	d := Discipline{
		Name:            r.FormValue("name"),
		NoOfTeamPlayers: no_of_team_players,
		TournamentID:    tournamentID,
	}

	if errs := validation.ValidateStruct(d); len(errs) > 0 {
		disciplines, err := h.service.GetAll(tournamentID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.RenderTemplate(w, r, h.templates["disciplines"], map[string]interface{}{
			"Errors":       errs,
			"Disciplines":  disciplines,
			"TournamentID": tournamentID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	_, err = h.service.Add(d)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID+"/disciplines", http.StatusSeeOther)
}

func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")

	disciplines, err := h.service.GetAll(tournamentID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["disciplines"], map[string]interface{}{
		"Disciplines":  disciplines,
		"TournamentID": tournamentID,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ShowHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")
	disciplineID := r.PathValue("discipline_id")

	discipline, err := h.service.Show(tournamentID, disciplineID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["discipline"], map[string]interface{}{
		"Discipline":   discipline,
		"TournamentID": tournamentID,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")
	disciplineID := r.PathValue("discipline_id")

	no_of_team_players, err := strconv.Atoi(r.FormValue("no_of_team_players"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	d := Discipline{
		Name:            r.FormValue("name"),
		TournamentID:    tournamentID,
		NoOfTeamPlayers: no_of_team_players,
		ID:              disciplineID,
	}

	if errs := validation.ValidateStruct(d); len(errs) > 0 {
		d, err := h.service.Show(tournamentID, disciplineID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.RenderTemplate(w, r, h.templates["discipline"], map[string]interface{}{
			"Errors":       errs,
			"Discipline":   d,
			"TournamentID": tournamentID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	_, err = h.service.Update(d)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID+"/disciplines/"+d.ID, http.StatusSeeOther)
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.PathValue("tournament_id")
	disciplineID := r.PathValue("discipline_id")

	err := h.service.Delete(tournamentID, disciplineID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/tournaments/"+tournamentID+"/disciplines", http.StatusSeeOther)
}
