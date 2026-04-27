package attendant

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tournament_manager/internal/player"
	"tournament_manager/internal/testutil"
)

const (
	testTournamentID = "11111111-1111-4111-8111-111111111111"
	testAttendantID  = "22222222-2222-4222-8222-222222222222"
)

type mockService struct {
	AddFunc    func(a Attendant) (Attendant, error)
	GetAllFunc func(tournamentID string) ([]DisplayAttendant, error)
	ShowFunc   func(tournamentID, id string) (DisplayAttendant, error)
	DeleteFunc func(tournamentID, id string) error
}

func (s *mockService) Add(a Attendant) (Attendant, error) {
	return s.AddFunc(a)
}

func (s *mockService) GetAll(tournamentID string) ([]DisplayAttendant, error) {
	return s.GetAllFunc(tournamentID)
}

func (s *mockService) Show(tournamentID, id string) (DisplayAttendant, error) {
	return s.ShowFunc(tournamentID, id)
}

func (s *mockService) Delete(tournamentID, id string) error {
	return s.DeleteFunc(tournamentID, id)
}

func newService() *mockService {
	return &mockService{
		AddFunc:    func(a Attendant) (Attendant, error) { return Attendant{}, nil },
		GetAllFunc: func(tournamentID string) ([]DisplayAttendant, error) { return []DisplayAttendant{}, nil },
		ShowFunc:   func(tournamentID, id string) (DisplayAttendant, error) { return DisplayAttendant{}, nil },
		DeleteFunc: func(tournamentID, id string) error { return nil },
	}
}

func fakeTemplates() map[string]*template.Template {
	t := template.Must(template.New("test").Parse(`{{define "base.html"}}OK{{end}}`))
	return map[string]*template.Template{
		"attendants": t,
		"attendant":  t,
	}
}

func newTestHandler(s ServiceInterface) *Handler {
	return &Handler{
		service:   s,
		templates: fakeTemplates(),
	}
}

func TestShowHandler_Success(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	called := false
	service.ShowFunc = func(tournamentID, id string) (DisplayAttendant, error) {
		called = true
		return DisplayAttendant{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/attendants"+testAttendantID, nil)
	w := httptest.NewRecorder()

	handler.ShowHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected OK, got %v", w.Code)
	}

	if !called {
		t.Fatal("expted Show to be called")
	}
}

func TestShowHandler_Failure(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	service.ShowFunc = func(tournamentID, id string) (DisplayAttendant, error) {
		return DisplayAttendant{}, errors.New("Error")
	}

	req := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/attendants"+testAttendantID, nil)
	w := httptest.NewRecorder()

	handler.ShowHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %v", w.Code)
	}
}

func TestListHandler_Success(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	called := false
	service.GetAllFunc = func(tournamentID string) ([]DisplayAttendant, error) {
		called = true
		return []DisplayAttendant{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/attendants", nil)
	w := httptest.NewRecorder()

	playerService := setupPlayerService(t)

	handler.ListHandler(w, req, playerService)

	if w.Code != http.StatusOK {
		t.Fatalf("expected OK, got %v", w.Code)
	}

	if !called {
		t.Fatal("expted Show to be called")
	}
}

func TestListHandler_Failure(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	service.GetAllFunc = func(tournamentID string) ([]DisplayAttendant, error) {
		return []DisplayAttendant{}, errors.New("Error")
	}

	req := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID+"/attendants", nil)
	w := httptest.NewRecorder()

	playerService := setupPlayerService(t)

	handler.ListHandler(w, req, playerService)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %v", w.Code)
	}
}

func TestCreateHandler_Success(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	called := false
	service.AddFunc = func(a Attendant) (Attendant, error) {
		called = true
		return Attendant{}, nil
	}

	form := strings.NewReader("player_id=123")
	req := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/attendants", form)
	w := httptest.NewRecorder()

	handler.CreateHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %v", w.Code)
	}

	if !called {
		t.Fatal("expted Show to be called")
	}
}

func TestCreateHandler_Failure(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	service.AddFunc = func(a Attendant) (Attendant, error) {
		return Attendant{}, errors.New("Error")
	}

	form := strings.NewReader("player_id=123")
	req := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/attendants", form)
	w := httptest.NewRecorder()

	handler.CreateHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v", w.Code)
	}
}

func TestDeleteHandler_Success(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	called := false
	service.DeleteFunc = func(tournamentID, id string) error {
		called = true
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID+"/attendants"+testAttendantID, nil)
	w := httptest.NewRecorder()

	handler.DeleteHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %v", w.Code)
	}

	if !called {
		t.Fatal("expted Show to be called")
	}
}

func TestDeleteHandler_Failure(t *testing.T) {
	service := newService()
	handler := newTestHandler(service)

	service.DeleteFunc = func(tournamentID, id string) error {
		return errors.New("Error")
	}

	req := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID+"/attendants"+testAttendantID, nil)
	w := httptest.NewRecorder()

	handler.DeleteHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v", w.Code)
	}
}

func setupPlayerService(t *testing.T) *player.Service {
	dbConn := testutil.SetupTestRepo(t)
	playerRepo := player.NewRepository(dbConn)
	return player.NewService(playerRepo)
}
