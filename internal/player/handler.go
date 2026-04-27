package player

import (
	"html/template"
	"net/http"
	"strings"
	tmpl "tournament_manager/internal/tmpl"
	"tournament_manager/internal/validation"

	"github.com/google/uuid"
)

type Handler struct {
	service   ServiceInterface
	templates map[string]*template.Template
}

func NewHandler(s ServiceInterface) *Handler {
	t := make(map[string]*template.Template)
	t["players"] = template.Must(template.New("").
		Funcs(tmpl.FuncMap()).
		ParseFiles(
			"templates/base.html",
			"templates/players.html",
		))

	return &Handler{
		service:   s,
		templates: t,
	}
}

func (h *Handler) PlayersHandler(w http.ResponseWriter, r *http.Request) {
	playerID, hasID, ok := parsePlayerPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if hasID {
		if _, err := uuid.Parse(playerID); err != nil {
			http.Error(w, "Invalid player ID", http.StatusNotFound)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		// TODO: miss show, do we need it?
		if hasID {
			http.NotFound(w, r)
			return
		}
		h.listHandler(w, r)
	case http.MethodPost:
		if hasID {
			http.NotFound(w, r)
			return
		}
		h.createHandler(w, r)
	case http.MethodPatch, http.MethodPut:
		if !hasID {
			http.NotFound(w, r)
			return
		}
		h.updateHandler(w, r, playerID)
	case http.MethodDelete:
		if !hasID {
			http.NotFound(w, r)
			return
		}
		h.deleteHandler(w, r, playerID)
	default:
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listHandler(w http.ResponseWriter, r *http.Request) {
	h.renderPlayers(w, r, Player{}, nil)
}

func (h *Handler) createHandler(w http.ResponseWriter, r *http.Request) {
	p, errs := playerFromRequest(r, "")
	if len(errs) > 0 {
		h.renderPlayers(w, r, p, errs)
		return
	}

	if _, err := h.service.Create(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/players", http.StatusSeeOther)
}

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request, playerID string) {
	p, errs := playerFromRequest(r, playerID)
	if len(errs) > 0 {
		h.renderPlayers(w, r, p, errs)
		return
	}

	if err := h.service.Update(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/players", http.StatusSeeOther)
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request, playerID string) {
	if err := h.service.Delete(playerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/players", http.StatusSeeOther)
}

func (h *Handler) renderPlayers(w http.ResponseWriter, r *http.Request, p Player, errs map[string]string) {
	players, err := h.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["players"], map[string]interface{}{
		"Errors":  errs,
		"Player":  p,
		"Players": players,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func playerFromRequest(r *http.Request, playerID string) (Player, map[string]string) {
	p := Player{
		ID:        playerID,
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Gender:    r.FormValue("gender"),
	}

	errs := map[string]string{}
	for field, msg := range validation.ValidateStruct(p) {
		errs[field] = msg
	}

	return p, errs
}

func parsePlayerPath(path string) (string, bool, bool) {
	if path == "/players" {
		return "", false, true
	}

	rest := strings.TrimPrefix(path, "/players/")
	if rest == path || rest == "" || strings.Contains(rest, "/") {
		return "", false, false
	}

	return rest, true, true
}
