package player

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tournament_manager/internal/tournament"
)

type mockService struct {
	CreateFunc func(Player) (Player, error)
	ListFunc   func(string) ([]Player, error)
	ShowFunc   func(string) (Player, error)
	UpdateFunc func(Player) error
	DeleteFunc func(string) error
}

func (m *mockService) Create(p Player) (Player, error) {
	return m.CreateFunc(p)
}

func (m *mockService) List(tournamentID string) ([]Player, error) {
	return m.ListFunc(tournamentID)
}

func (m *mockService) Show(id string) (Player, error) {
	return m.ShowFunc(id)
}

func (m *mockService) Update(p Player) error {
	return m.UpdateFunc(p)
}

func (m *mockService) Delete(id string) error {
	return m.DeleteFunc(id)
}

type mockTournamentService struct {
	ShowFunc func(string) (tournament.Tournament, error)
}

func (m *mockTournamentService) Show(id string) (tournament.Tournament, error) {
	return m.ShowFunc(id)
}

func newMockService() *mockService {
	return &mockService{
		CreateFunc: func(Player) (Player, error) { return Player{}, nil },
		ListFunc:   func(string) ([]Player, error) { return []Player{}, nil },
		ShowFunc:   func(string) (Player, error) { return Player{}, nil },
		UpdateFunc: func(Player) error { return nil },
		DeleteFunc: func(string) error { return nil },
	}
}

func newMockTournamentService() *mockTournamentService {
	return &mockTournamentService{
		ShowFunc: func(id string) (tournament.Tournament, error) {
			return tournament.Tournament{ID: id, Name: "Tournament", Location: "Berlin"}, nil
		},
	}
}

func fakeTemplates() map[string]*template.Template {
	t := template.Must(template.New("test").Parse(`{{define "base.html"}}OK{{end}}`))
	return map[string]*template.Template{
		"tournament": t,
	}
}

func newTestHandler(s ServiceInterface, ts TournamentService) *Handler {
	return &Handler{
		service:           s,
		tournamentService: ts,
		templates:         fakeTemplates(),
	}
}

func TestCreateHandler_Success(t *testing.T) {
	mock := newMockService()
	called := false
	mock.CreateFunc = func(p Player) (Player, error) {
		called = true
		return Player{}, nil
	}

	h := newTestHandler(mock, newMockTournamentService())
	req := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/players", validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", w.Code)
	}

	if got := w.Header().Get("Location"); got != "/tournaments/"+testTournamentID {
		t.Fatalf("expected redirect to tournament, got %s", got)
	}

	if !called {
		t.Fatalf("expected Create to be called")
	}
}

func TestCreateHandler_ValidationFail(t *testing.T) {
	h := newTestHandler(newMockService(), newMockTournamentService())
	req := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/players", strings.NewReader("first_name=&last_name=&gender=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateHandler_InvalidTournamentID(t *testing.T) {
	h := newTestHandler(newMockService(), newMockTournamentService())
	req := httptest.NewRequest(http.MethodPost, "/tournaments/not-a-uuid/players", validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateHandler_ServiceFail(t *testing.T) {
	mock := newMockService()
	mock.CreateFunc = func(p Player) (Player, error) {
		return Player{}, errors.New("errors")
	}

	h := newTestHandler(mock, newMockTournamentService())
	req := httptest.NewRequest(http.MethodPost, "/tournaments/"+testTournamentID+"/players", validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateHandler_Success(t *testing.T) {
	mock := newMockService()
	called := false
	mock.UpdateFunc = func(p Player) error {
		called = true
		return nil
	}

	h := newTestHandler(mock, newMockTournamentService())
	req := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID+"/players/"+testPlayerID, validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expected Update to be called")
	}
}

func TestUpdateHandler_ValidationFail(t *testing.T) {
	h := newTestHandler(newMockService(), newMockTournamentService())
	req := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID+"/players/"+testPlayerID, strings.NewReader("first_name=&last_name=&gender=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateHandler_InvalidPlayerID(t *testing.T) {
	h := newTestHandler(newMockService(), newMockTournamentService())
	req := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID+"/players/not-a-uuid", validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateHandler_ServiceFail(t *testing.T) {
	mock := newMockService()
	mock.UpdateFunc = func(p Player) error {
		return errors.New("errors")
	}

	h := newTestHandler(mock, newMockTournamentService())
	req := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID+"/players/"+testPlayerID, validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteHandler_Success(t *testing.T) {
	mock := newMockService()
	called := false
	mock.DeleteFunc = func(id string) error {
		called = true
		return nil
	}

	h := newTestHandler(mock, newMockTournamentService())
	req := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID+"/players/"+testPlayerID, nil)
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expected Delete to be called")
	}
}

func TestDeleteHandler_ServiceFail(t *testing.T) {
	mock := newMockService()
	mock.DeleteFunc = func(id string) error {
		return errors.New("errors")
	}

	h := newTestHandler(mock, newMockTournamentService())
	req := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID+"/players/"+testPlayerID, nil)
	w := httptest.NewRecorder()

	h.ByTournamentHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func validPlayerForm() *strings.Reader {
	return strings.NewReader("first_name=Foo&last_name=Bar&gender=male")
}
