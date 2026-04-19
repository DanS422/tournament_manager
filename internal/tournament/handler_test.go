package tournament

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockService struct {
	ListFunc   func() ([]Tournament, error)
	CreateFunc func(string, string) (Tournament, error)
	ShowFunc   func(string) (Tournament, error)
	UpdateFunc func(string, string, string) error
	DeleteFunc func(string) error
}

func (m *mockService) List() ([]Tournament, error) {
	return m.ListFunc()
}

func (m *mockService) Create(name, location string) (Tournament, error) {
	return m.CreateFunc(name, location)
}

func (m *mockService) Show(id string) (Tournament, error) {
	return m.ShowFunc(id)
}

func (m *mockService) Update(id string, name, location string) error {
	return m.UpdateFunc(id, name, location)
}

func (m *mockService) Delete(id string) error {
	return m.DeleteFunc(id)
}

func fakeTemplates() map[string]*template.Template {
	t := template.Must(template.New("test").Parse(`{{define "base.html"}}OK{{end}}`))
	return map[string]*template.Template{
		"tournaments": t,
		"tournament":  t,
	}
}

func brokenTemplates() map[string]*template.Template {
	t := template.Must(template.New("test").Parse(`{{define "base.html"}}{{template "missing" .}}{{end}}`))
	return map[string]*template.Template{
		"tournaments": t,
		"tournament":  t,
	}
}

func newMockService() *mockService {
	return &mockService{
		ListFunc:   func() ([]Tournament, error) { return nil, nil },
		CreateFunc: func(string, string) (Tournament, error) { return Tournament{}, nil },
		ShowFunc:   func(string) (Tournament, error) { return Tournament{}, nil },
		UpdateFunc: func(string, string, string) error { return nil },
		DeleteFunc: func(string) error { return nil },
	}
}

func TestCreateHandler_ValidationFail(t *testing.T) {
	mock := newMockService()

	mock.ListFunc = func() ([]Tournament, error) {
		return []Tournament{}, nil
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	form := strings.NewReader("name=&location=")
	req := httptest.NewRequest(http.MethodPost, "/", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.CreateHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateHandler_Success(t *testing.T) {
	mock := newMockService()

	called := false
	mock.CreateFunc = func(name, location string) (Tournament, error) {
		called = true
		return Tournament{}, nil
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	form := strings.NewReader("name=Test&location=SG")
	req := httptest.NewRequest(http.MethodPost, "/", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.CreateHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expected Create to be called")
	}
}

func TestListHandler_Success(t *testing.T) {
	mock := newMockService()
	called := false
	mock.ListFunc = func() ([]Tournament, error) {
		called = true
		return []Tournament{}, nil
	}
	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/html")

	w := httptest.NewRecorder()

	h.ListHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected OK, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expected List to be called")
	}
}

func TestListHandler_Fail(t *testing.T) {
	mock := newMockService()
	called := false
	mock.ListFunc = func() ([]Tournament, error) {
		called = true
		return []Tournament{}, errors.New("errors")
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/html")

	w := httptest.NewRecorder()

	h.ListHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expected List to be called")
	}
}

func TestListHandler_RenderFail(t *testing.T) {
	mock := newMockService()
	mock.ListFunc = func() ([]Tournament, error) {
		return []Tournament{}, nil
	}

	h := &Handler{
		service:   mock,
		templates: brokenTemplates(),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ListHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestListHandler_UnknownPath(t *testing.T) {
	mock := newMockService()
	called := false
	mock.ListFunc = func() ([]Tournament, error) {
		called = true
		return []Tournament{}, nil
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	req := httptest.NewRequest(http.MethodGet, "/not-a-real-page", nil)
	w := httptest.NewRecorder()

	h.ListHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	if called {
		t.Fatalf("expected List not to be called")
	}
}

func TestByIDHandler_Fail(t *testing.T) {
	mock := newMockService()
	mock.ShowFunc = func(id string) (Tournament, error) {
		return Tournament{}, errors.New("errors")
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	req := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID, nil)
	req.Header.Set("Content-Type", "application/html")
	w := httptest.NewRecorder()
	h.ByIDHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expect 404, got %d", w.Code)
	}
}

func TestByIDHandler_Show_Sucess(t *testing.T) {
	mock := newMockService()
	called := false
	mock.ShowFunc = func(id string) (Tournament, error) {
		called = true
		return Tournament{}, nil
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	req := httptest.NewRequest(http.MethodGet, "/tournaments/"+testTournamentID, nil)
	req.Header.Set("Content-Type", "application/html")
	w := httptest.NewRecorder()
	h.ByIDHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expect Show to be called")
	}
}

func TestByIDHandler_Delete_Success(t *testing.T) {
	mock := newMockService()
	called := false
	mock.DeleteFunc = func(id string) error {
		called = true
		return nil
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	req := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID, nil)
	req.Header.Set("Content-Type", "application/html")

	w := httptest.NewRecorder()
	h.ByIDHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expect 302, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expect Delete to be called")
	}
}

func TestByIDHandler_Delete_Fail(t *testing.T) {
	mock := newMockService()
	called := false
	mock.DeleteFunc = func(id string) error {
		called = true
		return errors.New("errors")
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	req := httptest.NewRequest(http.MethodDelete, "/tournaments/"+testTournamentID, nil)
	req.Header.Set("Content-Type", "application/html")

	w := httptest.NewRecorder()
	h.ByIDHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expect 400, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expect Delete to be called")
	}
}

func TestByIDHandler_Update_Sucess(t *testing.T) {
	mock := newMockService()
	called := false
	mock.UpdateFunc = func(id string, name string, location string) error {
		called = true
		return nil
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	form := strings.NewReader("name=Test&location=SG")
	req := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID, form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	h.ByIDHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expect 302, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expect Update to be called")
	}
}

func TestByIDHandler_Update_Fail(t *testing.T) {
	mock := newMockService()
	called := false
	mock.UpdateFunc = func(id string, name string, location string) error {
		called = true
		return errors.New("error")
	}

	h := &Handler{
		service:   mock,
		templates: fakeTemplates(),
	}

	form := strings.NewReader("name=Test&location=SG")
	req := httptest.NewRequest(http.MethodPatch, "/tournaments/"+testTournamentID, form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	h.ByIDHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expect 400, got %d", w.Code)
	}

	if !called {
		t.Fatalf("expect Update to be called")
	}
}
