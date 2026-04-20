package player

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockService struct {
	CreateFunc func(Player) (Player, error)
	ListFunc   func() ([]Player, error)
	ShowFunc   func(string) (Player, error)
	UpdateFunc func(Player) error
	DeleteFunc func(string) error
}

func (m *mockService) Create(p Player) (Player, error) {
	return m.CreateFunc(p)
}

func (m *mockService) List() ([]Player, error) {
	return m.ListFunc()
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

func newMockService() *mockService {
	return &mockService{
		CreateFunc: func(Player) (Player, error) { return Player{}, nil },
		ListFunc:   func() ([]Player, error) { return []Player{}, nil },
		ShowFunc:   func(string) (Player, error) { return Player{}, nil },
		UpdateFunc: func(Player) error { return nil },
		DeleteFunc: func(string) error { return nil },
	}
}

func fakeTemplates() map[string]*template.Template {
	t := template.Must(template.New("test").Parse(`{{define "base.html"}}OK{{end}}`))
	return map[string]*template.Template{
		"players": t,
	}
}

func newTestHandler(s ServiceInterface) *Handler {
	return &Handler{
		service:   s,
		templates: fakeTemplates(),
	}
}

func TestListHandler_Success(t *testing.T) {
	mock := newMockService()
	called := false
	mock.ListFunc = func() ([]Player, error) {
		called = true
		return []Player{}, nil
	}

	h := newTestHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/players", nil)
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expected List to be called")
	}
}

func TestCreateHandler_Success(t *testing.T) {
	mock := newMockService()
	called := false
	mock.CreateFunc = func(p Player) (Player, error) {
		called = true
		return Player{}, nil
	}

	h := newTestHandler(mock)
	req := httptest.NewRequest(http.MethodPost, "/players", validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", w.Code)
	}

	if got := w.Header().Get("Location"); got != "/players" {
		t.Fatalf("expected redirect to players, got %s", got)
	}

	if !called {
		t.Fatalf("expected Create to be called")
	}
}

func TestCreateHandler_ValidationFail(t *testing.T) {
	h := newTestHandler(newMockService())
	req := httptest.NewRequest(http.MethodPost, "/players", strings.NewReader("first_name=&last_name=&gender=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateHandler_ServiceFail(t *testing.T) {
	mock := newMockService()
	mock.CreateFunc = func(p Player) (Player, error) {
		return Player{}, errors.New("errors")
	}

	h := newTestHandler(mock)
	req := httptest.NewRequest(http.MethodPost, "/players", validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

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

	h := newTestHandler(mock)
	req := httptest.NewRequest(http.MethodPatch, "/players/"+testPlayerID, validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expected Update to be called")
	}
}

func TestUpdateHandler_ValidationFail(t *testing.T) {
	h := newTestHandler(newMockService())
	req := httptest.NewRequest(http.MethodPatch, "/players/"+testPlayerID, strings.NewReader("first_name=&last_name=&gender=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateHandler_InvalidPlayerID(t *testing.T) {
	h := newTestHandler(newMockService())
	req := httptest.NewRequest(http.MethodPatch, "/players/not-a-uuid", validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateHandler_ServiceFail(t *testing.T) {
	mock := newMockService()
	mock.UpdateFunc = func(p Player) error {
		return errors.New("errors")
	}

	h := newTestHandler(mock)
	req := httptest.NewRequest(http.MethodPatch, "/players/"+testPlayerID, validPlayerForm())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

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

	h := newTestHandler(mock)
	req := httptest.NewRequest(http.MethodDelete, "/players/"+testPlayerID, nil)
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

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

	h := newTestHandler(mock)
	req := httptest.NewRequest(http.MethodDelete, "/players/"+testPlayerID, nil)
	w := httptest.NewRecorder()

	h.PlayersHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func validPlayerForm() *strings.Reader {
	return strings.NewReader("first_name=Foo&last_name=Bar&gender=male")
}
