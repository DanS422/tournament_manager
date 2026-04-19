package tournament

import (
	"errors"
	"testing"
)

type mockRepository struct {
	AddFunc    func(t Tournament) (Tournament, error)
	GetAllFunc func() ([]Tournament, error)
	ShowFunc   func(id string) (Tournament, error)
	UpdateFunc func(id string, t Tournament) error
	DeleteFunc func(id string) error
}

func (r *mockRepository) Add(t Tournament) (Tournament, error) {
	return r.AddFunc(t)
}

func (r *mockRepository) GetAll() ([]Tournament, error) {
	return r.GetAllFunc()
}

func (r *mockRepository) Show(id string) (Tournament, error) {
	return r.ShowFunc(id)
}

func (r *mockRepository) Update(id string, t Tournament) error {
	return r.UpdateFunc(id, t)
}

func (r *mockRepository) Delete(id string) error {
	return r.DeleteFunc(id)
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		AddFunc:    func(t Tournament) (Tournament, error) { return Tournament{}, nil },
		GetAllFunc: func() ([]Tournament, error) { return []Tournament{}, nil },
		ShowFunc:   func(id string) (Tournament, error) { return Tournament{}, nil },
		UpdateFunc: func(id string, t Tournament) error { return nil },
		DeleteFunc: func(id string) error { return nil },
	}
}

func TestCreate_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.AddFunc = func(t Tournament) (Tournament, error) {
		return Tournament{}, errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Create("foo", "bar")

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestCreate_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.AddFunc = func(t Tournament) (Tournament, error) {
		called = true
		return Tournament{}, nil
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Create("foo", "bar")

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Create to be called")
	}
}

func TestList_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.GetAllFunc = func() ([]Tournament, error) {
		return []Tournament{}, errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.List()

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestList_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.GetAllFunc = func() ([]Tournament, error) {
		called = true
		return []Tournament{}, nil
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.List()

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Create to be called")
	}
}

func TestShow_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.ShowFunc = func(id string) (Tournament, error) {
		return Tournament{}, errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Show(testTournamentID)

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestShow_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.ShowFunc = func(id string) (Tournament, error) {
		called = true
		return Tournament{}, nil
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Show(testTournamentID)

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Create to be called")
	}
}

func TestUpdate_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.UpdateFunc = func(id string, t Tournament) error {
		return errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	err := service.Update(testTournamentID, "foo", "bar")

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestUpdate_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.UpdateFunc = func(id string, t Tournament) error {
		called = true
		return nil
	}
	service := &Service{
		repo: mock,
	}

	err := service.Update(testTournamentID, "foo", "bar")

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Update to be called")
	}
}

func TestDelete_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.DeleteFunc = func(id string) error {
		return errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	err := service.Delete(testTournamentID)

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestDelete_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.DeleteFunc = func(id string) error {
		called = true
		return nil
	}
	service := &Service{
		repo: mock,
	}

	err := service.Delete(testTournamentID)

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Update to be called")
	}
}
