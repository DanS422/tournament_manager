package player

import (
	"errors"
	"testing"
)

const testTournamentID = "33333333-3333-4333-8333-333333333333"

type mockRepository struct {
	AddFunc    func(p Player) (Player, error)
	GetAllFunc func(tournamentID string) ([]Player, error)
	ShowFunc   func(id string) (Player, error)
	UpdateFunc func(p Player) error
	DeleteFunc func(id string) error
}

func (r *mockRepository) Add(p Player) (Player, error) {
	return r.AddFunc(p)
}

func (r *mockRepository) GetAll(tournamentID string) ([]Player, error) {
	return r.GetAllFunc(tournamentID)
}

func (r *mockRepository) Show(id string) (Player, error) {
	return r.ShowFunc(id)
}

func (r *mockRepository) Update(p Player) error {
	return r.UpdateFunc(p)
}

func (r *mockRepository) Delete(id string) error {
	return r.DeleteFunc(id)
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		AddFunc:    func(p Player) (Player, error) { return Player{}, nil },
		GetAllFunc: func(tournamentID string) ([]Player, error) { return []Player{}, nil },
		ShowFunc:   func(id string) (Player, error) { return Player{}, nil },
		UpdateFunc: func(p Player) error { return nil },
		DeleteFunc: func(id string) error { return nil },
	}
}

func TestCreate_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.AddFunc = func(p Player) (Player, error) {
		return Player{}, errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Create(Player{FirstName: FirstName, LastName: LastName})

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestCreate_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.AddFunc = func(p Player) (Player, error) {
		called = true
		return Player{}, nil
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Create(Player{FirstName: FirstName, LastName: LastName})

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Create to be called")
	}
}

func TestList_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.GetAllFunc = func(tournamentID string) ([]Player, error) {
		return []Player{}, errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.List(testTournamentID)

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestList_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.GetAllFunc = func(tournamentID string) ([]Player, error) {
		called = true
		return []Player{}, nil
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.List(testTournamentID)

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected List to be called")
	}
}

func TestShow_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.ShowFunc = func(id string) (Player, error) {
		return Player{}, errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Show(testPlayerID)

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestShow_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.ShowFunc = func(id string) (Player, error) {
		called = true
		return Player{}, nil
	}

	service := &Service{
		repo: mock,
	}

	_, err := service.Show(testPlayerID)

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Show to be called")
	}
}

func TestUpdate_Fail(t *testing.T) {
	mock := newMockRepository()

	mock.UpdateFunc = func(p Player) error {
		return errors.New("errors")
	}

	service := &Service{
		repo: mock,
	}

	err := service.Update(Player{ID: testPlayerID, FirstName: FirstName, LastName: LastName})

	if err == nil {
		t.Fatalf("expect to error out")
	}
}

func TestUpdate_Success(t *testing.T) {
	mock := newMockRepository()

	called := false
	mock.UpdateFunc = func(p Player) error {
		called = true
		return nil
	}
	service := &Service{
		repo: mock,
	}

	err := service.Update(Player{ID: testPlayerID, FirstName: FirstName, LastName: LastName})

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

	err := service.Delete(testPlayerID)

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

	err := service.Delete(testPlayerID)

	if err != nil {
		t.Fatalf("expect no errors")
	}

	if !called {
		t.Fatalf("expected Delete to be called")
	}
}
