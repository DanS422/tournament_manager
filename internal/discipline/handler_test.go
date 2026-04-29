package discipline

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testTournamentID = "11111111-1111-4111-8111-111111111111"
	testDisciplineID = "22222222-2222-4222-8222-222222222222"
)

type mockService struct {
	AddFunc    func(d Discipline) (Discipline, error)
	GetAllFunc func(tournamentID string) ([]Discipline, error)
	ShowFunc   func(tournamentID, id string) (Discipline, error)
	UpdateFunc func(d Discipline) (Discipline, error)
	DeleteFunc func(tournamentID, id string) error
}

func (s *mockService) Add(d Discipline) (Discipline, error) {
	return s.AddFunc(d)
}

func (s *mockService) GetAll(tournamentID string) ([]Discipline, error) {
	return s.GetAllFunc(tournamentID)
}

func (s *mockService) Show(tournamentID, id string) (Discipline, error) {
	return s.ShowFunc(tournamentID, id)
}

func (s *mockService) Update(d Discipline) (Discipline, error) {
	return s.UpdateFunc(d)
}

func (s *mockService) Delete(tournamentID, id string) error {
	return s.DeleteFunc(tournamentID, id)
}

func mockTemplates() map[string]*template.Template {
	t := template.Must(template.New("test").Parse(`{{define "base.html"}}OK{{end}}`))

	return map[string]*template.Template{
		"disciplines": t,
		"discipline":  t,
	}
}

func newService() *mockService {
	return &mockService{
		AddFunc:    func(d Discipline) (Discipline, error) { return Discipline{}, nil },
		GetAllFunc: func(tournamentID string) ([]Discipline, error) { return []Discipline{}, nil },
		ShowFunc:   func(tournamentID, id string) (Discipline, error) { return Discipline{}, nil },
		UpdateFunc: func(d Discipline) (Discipline, error) { return Discipline{}, nil },
		DeleteFunc: func(tournamentID, id string) error { return nil },
	}
}

func newHandler(s ServiceInterface) *Handler {
	return &Handler{
		service:   s,
		templates: mockTemplates(),
	}
}

func TestHandlerCreate_Success(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	called := false
	service.AddFunc = func(d Discipline) (Discipline, error) {
		called = true
		return Discipline{}, nil
	}

	requestForm := strings.NewReader("name=foo&no_of_team_players=1")
	r := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/disciplines", requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)

	w := httptest.NewRecorder()

	handler.CreateHandler(w, r)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("Expected 303, got %v", w.Code)
	}

	if !called {
		t.Fatal("Expected Add function to be called")
	}
}

func TestHandlerCreate_Failure(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	service.AddFunc = func(d Discipline) (Discipline, error) {
		return Discipline{}, errors.New("Error")
	}

	requestForm := strings.NewReader("name=foo&no_of_team_players=1")
	r := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/disciplines", requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)

	w := httptest.NewRecorder()

	handler.CreateHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %v", w.Code)
	}
}

func TestHandlerCreate_ValidationFailure(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	service.AddFunc = func(d Discipline) (Discipline, error) {
		return Discipline{}, errors.New("Error")
	}

	requestForm := strings.NewReader("name=")
	r := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/disciplines", requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)

	w := httptest.NewRecorder()

	handler.CreateHandler(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Expected 422, got %v", w.Code)
	}
}

func TestHandlerList_Success(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	called := false
	service.GetAllFunc = func(tournamentID string) ([]Discipline, error) {
		called = true
		return []Discipline{}, nil
	}

	r := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/disciplines", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)

	w := httptest.NewRecorder()

	handler.ListHandler(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %v", w.Code)
	}

	if !called {
		t.Fatal("Expected GetAll function to be called")
	}
}

func TestHandlerList_Failure(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	service.GetAllFunc = func(tournamentID string) ([]Discipline, error) {
		return []Discipline{}, errors.New("Errors")
	}

	r := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/disciplines", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)

	w := httptest.NewRecorder()

	handler.ListHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %v", w.Code)
	}
}

func TestHandlerShow_Success(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	called := false
	service.ShowFunc = func(tournamentID, id string) (Discipline, error) {
		called = true
		return Discipline{}, nil
	}

	r := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/disciplines/"+testDisciplineID, nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)
	r.SetPathValue("discipline_id", testDisciplineID)

	w := httptest.NewRecorder()

	handler.ShowHandler(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %v", w.Code)
	}

	if !called {
		t.Fatal("Expected Show function to be called")
	}
}

func TestHandlerShow_Failure(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	service.ShowFunc = func(tournamentID, id string) (Discipline, error) {
		return Discipline{}, errors.New("errors")
	}

	r := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/disciplines/"+testDisciplineID, nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)
	r.SetPathValue("discipline_id", testDisciplineID)

	w := httptest.NewRecorder()

	handler.ShowHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 500, got %v", w.Code)
	}
}

func TestHandlerUpdate_Success(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	called := false
	service.UpdateFunc = func(d Discipline) (Discipline, error) {
		called = true
		return Discipline{}, nil
	}

	requestForm := strings.NewReader("name=foo&no_of_team_players=2")
	r := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID+"/disciplines/"+testDisciplineID, requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)
	r.SetPathValue("discipline_id", testDisciplineID)

	w := httptest.NewRecorder()

	handler.UpdateHandler(w, r)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("Expected 303, got %v", w.Code)
	}

	if !called {
		t.Fatal("Expected Update function to be called")
	}
}

func TestHandlerUpdate_Failure(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	service.UpdateFunc = func(d Discipline) (Discipline, error) {
		return Discipline{}, errors.New("Error")
	}

	requestForm := strings.NewReader("name=foo&no_of_team_players=1")
	r := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID+"/disciplines/"+testDisciplineID, requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)
	r.SetPathValue("discipline_id", testDisciplineID)

	w := httptest.NewRecorder()

	handler.UpdateHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected 400, got %v", w.Code)
	}
}

func TestHandlerUpdate_ValidationFailure(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	service.UpdateFunc = func(d Discipline) (Discipline, error) {
		return Discipline{}, errors.New("Error")
	}

	requestForm := strings.NewReader("name=")
	r := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID+"/disciplines/"+testDisciplineID, requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)
	r.SetPathValue("discipline_id", testDisciplineID)

	w := httptest.NewRecorder()

	handler.UpdateHandler(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Expected 422, got %v", w.Code)
	}
}

func TestHandlerDelete_Success(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	called := false
	service.DeleteFunc = func(tournamentID, id string) error {
		called = true
		return nil
	}

	requestForm := strings.NewReader("name=foo")
	r := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID+"/disciplines/"+testDisciplineID, requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)
	r.SetPathValue("discipline_id", testDisciplineID)

	w := httptest.NewRecorder()

	handler.DeleteHandler(w, r)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("Expected 303, got %v", w.Code)
	}

	if !called {
		t.Fatal("Expected Update function to be called")
	}
}

func TestHandlerDelete_Failure(t *testing.T) {
	service := newService()
	handler := newHandler(service)

	service.DeleteFunc = func(tournamentID, id string) error {
		return errors.New("Errors")
	}

	requestForm := strings.NewReader("name=foo")
	r := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID+"/disciplines/"+testDisciplineID, requestForm)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("tournament_id", testTournamentID)
	r.SetPathValue("discipline_id", testDisciplineID)

	w := httptest.NewRecorder()

	handler.DeleteHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %v", w.Code)
	}
}
