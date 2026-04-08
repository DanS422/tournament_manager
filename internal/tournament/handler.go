package tournament

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
	tmpl    *template.Template
}

func NewHandler(s *Service) *Handler {
	tmpl := template.Must(template.ParseFiles("templates/tournament.html", "templates/tournaments.html"))
	return &Handler{service: s, tmpl: tmpl}
}

func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	tournaments := h.service.List()
	h.tmpl.ExecuteTemplate(w, "tournaments.html", tournaments)
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	name := r.FormValue("name")
	location := r.FormValue("location")
	_, err := h.service.Create(name, location)
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
		h.showHandler(w, id)
	default:
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) showHandler(w http.ResponseWriter, id int) {
	tournament, err := h.service.show(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	h.tmpl.ExecuteTemplate(w, "tournament.html", tournament)
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
