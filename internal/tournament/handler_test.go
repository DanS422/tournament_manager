package tournament

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockService struct {
	ListFunc   func() ([]Tournament, error)
	CreateFunc func(string, string) (Tournament, error)
	ShowFunc   func(int) (Tournament, error)
	UpdateFunc func(int, string, string) error
	DeleteFunc func(int) error
}

func (m *mockService) List() ([]Tournament, error) {
	return m.ListFunc()
}

func (m *mockService) Create(name, location string) (Tournament, error) {
	return m.CreateFunc(name, location)
}

func (m *mockService) Show(id int) (Tournament, error) {
	return m.ShowFunc(id)
}

func (m *mockService) Update(id int, name, location string) error {
	return m.UpdateFunc(id, name, location)
}

func (m *mockService) Delete(id int) error {
	return m.DeleteFunc(id)
}

func fakeTemplates() map[string]*template.Template {
	t := template.Must(template.New("test").Parse("OK"))
	return map[string]*template.Template{
		"tournaments": t,
		"tournament":  t,
	}
}

func newMockService() *mockService {
	return &mockService{
		ListFunc:   func() ([]Tournament, error) { return nil, nil },
		CreateFunc: func(string, string) (Tournament, error) { return Tournament{}, nil },
		ShowFunc:   func(int) (Tournament, error) { return Tournament{}, nil },
		UpdateFunc: func(int, string, string) error { return nil },
		DeleteFunc: func(int) error { return nil },
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
