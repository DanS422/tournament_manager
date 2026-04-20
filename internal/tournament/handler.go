package tournament

import (
	"html/template"
	"net/http"
	"strings"
	tmpl "tournament_manager/internal/tmpl"
	"tournament_manager/internal/validation"

	"github.com/google/uuid"
)

type Handler struct {
	service    ServiceInterface
	templates  map[string]*template.Template
	playerList func(tournamentID string) (any, error)
}

func NewHandler(s ServiceInterface) *Handler {
	t := make(map[string]*template.Template)

	pages := []string{"tournaments", "tournament"}
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

func (h *Handler) SetPlayerListFunc(fn func(tournamentID string) (any, error)) {
	h.playerList = fn
}

func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tournaments, err := h.service.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["tournaments"], map[string]interface{}{
		"Tournaments": tournaments,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t := Tournament{
		Name:     r.FormValue("name"),
		Location: r.FormValue("location"),
	}
	if errs := validation.ValidateStruct(t); len(errs) > 0 {
		tournaments, err := h.service.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.RenderTemplate(w, r, h.templates["tournaments"], map[string]interface{}{
			"Errors":      errs,
			"Tournaments": tournaments,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
	_, err := h.service.Create(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tournaments/")

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	switch requestMethod := r.Method; requestMethod {
	case http.MethodPatch, http.MethodPut:
		h.updateHandler(w, r, id)
	case http.MethodDelete:
		h.deleteHandler(w, r, id)
	case http.MethodGet:
		h.showHandler(w, r, id)
	default:
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) showHandler(w http.ResponseWriter, r *http.Request, id string) {
	tournament, err := h.service.Show(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	players, err := h.playersFor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.RenderTemplate(w, r, h.templates["tournament"], map[string]interface{}{
		"Player":     map[string]string{},
		"Tournament": tournament,
		"Players":    players,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request, id string) {
	t := Tournament{
		ID:       id,
		Name:     r.FormValue("name"),
		Location: r.FormValue("location"),
	}
	if errs := validation.ValidateStruct(t); len(errs) > 0 {
		players, err := h.playersFor(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.RenderTemplate(w, r, h.templates["tournament"], map[string]interface{}{
			"Errors":     errs,
			"Player":     map[string]string{},
			"Players":    players,
			"Tournament": t,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	err := h.service.Update(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) playersFor(tournamentID string) (any, error) {
	if h.playerList == nil {
		return []any{}, nil
	}

	return h.playerList(tournamentID)
}
