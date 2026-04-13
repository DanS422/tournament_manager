package tournament

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	tmpl "tournament_manager/internal/tmpl"
	"tournament_manager/internal/validation"
)

type Handler struct {
	service   *Service
	templates map[string]*template.Template
}

func NewHandler(s *Service) *Handler {
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

func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	tournaments, err := h.service.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.RenderTemplate(w, r, h.templates["tournaments"], map[string]interface{}{
		"Tournaments": tournaments,
	})
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
		fmt.Println(errs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tmpl.RenderTemplate(w, r, h.templates["tournaments"], map[string]interface{}{
			"Errors":      errs,
			"Tournaments": tournaments,
		})

		return
	}
	_, err := h.service.Create(t.Name, t.Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tournaments/")
	id, _ := strconv.Atoi(idStr)

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

func (h *Handler) showHandler(w http.ResponseWriter, r *http.Request, id int) {
	tournament, err := h.service.show(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	tmpl.RenderTemplate(w, r, h.templates["tournament"], map[string]interface{}{
		"Tournament": tournament,
	})
}

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request, id int) {
	name := r.FormValue("name")
	location := r.FormValue("location")

	err := h.service.Update(id, name, location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
