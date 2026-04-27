package attendant

import (
	"errors"
	"testing"
)

type mockRepository struct {
	AddFunc    func(a Attendant) (Attendant, error)
	GetAllFunc func(tournamentID string) ([]DisplayAttendant, error)
	ShowFunc   func(tournamentID, id string) (DisplayAttendant, error)
	DeleteFunc func(tournamentID, id string) error
}

func (r *mockRepository) Add(a Attendant) (Attendant, error) {
	return r.AddFunc(a)
}

func (r *mockRepository) GetAll(tournamentID string) ([]DisplayAttendant, error) {
	return r.GetAllFunc(tournamentID)
}

func (r *mockRepository) Show(tournamentID, id string) (DisplayAttendant, error) {
	return r.ShowFunc(tournamentID, id)
}

func (r *mockRepository) Delete(tournamentID, id string) error {
	return r.DeleteFunc(tournamentID, id)
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		AddFunc:    func(a Attendant) (Attendant, error) { return Attendant{}, nil },
		GetAllFunc: func(tournamentID string) ([]DisplayAttendant, error) { return []DisplayAttendant{}, nil },
		ShowFunc:   func(tournamentID, id string) (DisplayAttendant, error) { return DisplayAttendant{}, nil },
		DeleteFunc: func(tournamentID, id string) error { return nil },
	}
}

func TestServiceAdd_Success(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	called := false
	repo.AddFunc = func(a Attendant) (Attendant, error) {
		called = true
		return Attendant{}, nil
	}

	_, err := s.Add(Attendant{})

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if !called {
		t.Fatal("expected add function is being called")
	}
}

func TestServiceAdd_Failure(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	repo.AddFunc = func(a Attendant) (Attendant, error) {
		return Attendant{}, errors.New("Errors")
	}

	_, err := s.Add(Attendant{})

	if err == nil {
		t.Fatal("expected error to be raised")
	}
}

func TestServiceGetAll_Success(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	called := false
	repo.GetAllFunc = func(tournamentID string) ([]DisplayAttendant, error) {
		called = true
		return []DisplayAttendant{}, nil
	}

	_, err := s.GetAll("1234")

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if !called {
		t.Fatal("expected add function is being called")
	}
}

func TestServiceGetAll_Failure(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	repo.GetAllFunc = func(tournamentID string) ([]DisplayAttendant, error) {
		return []DisplayAttendant{}, errors.New("Errors")
	}

	_, err := s.GetAll("1234")

	if err == nil {
		t.Fatal("expected error to be raised")
	}
}

func TestServiceShow_Success(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	called := false
	repo.ShowFunc = func(tournamentID, id string) (DisplayAttendant, error) {
		called = true
		return DisplayAttendant{}, nil
	}

	_, err := s.Show("1234", "xyz")

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if !called {
		t.Fatal("expected add function is being called")
	}
}

func TestServiceShow_Failure(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	repo.ShowFunc = func(tournamentID, id string) (DisplayAttendant, error) {
		return DisplayAttendant{}, errors.New("Errors")
	}

	_, err := s.Show("1234", "xyz")

	if err == nil {
		t.Fatal("expected error to be raised")
	}
}

func TestServiceDelete_Success(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	called := false
	repo.DeleteFunc = func(tournamentID, id string) error {
		called = true
		return nil
	}

	err := s.Delete("1234", "xyz")

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if !called {
		t.Fatal("expected add function is being called")
	}
}

func TestServiceDelete_Failure(t *testing.T) {
	repo := NewMockRepository()
	s := NewService(repo)

	repo.DeleteFunc = func(tournamentID, id string) error {
		return errors.New("Errors")
	}

	err := s.Delete("1234", "xyz")

	if err == nil {
		t.Fatal("expected error to be raised")
	}
}
