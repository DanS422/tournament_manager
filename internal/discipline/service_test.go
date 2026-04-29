package discipline

import (
	"errors"
	"testing"
)

type mockRepository struct {
	AddFunc    func(d Discipline) (Discipline, error)
	GetAllFunc func(tournamentID string) ([]Discipline, error)
	ShowFunc   func(tournamentID, id string) (Discipline, error)
	UpdateFunc func(d Discipline) (Discipline, error)
	DeleteFunc func(tournamentID, id string) error
}

func (r *mockRepository) Add(d Discipline) (Discipline, error) {
	return r.AddFunc(d)
}

func (r *mockRepository) GetAll(tournamentID string) ([]Discipline, error) {
	return r.GetAllFunc(tournamentID)
}

func (r *mockRepository) Show(tournamentID, id string) (Discipline, error) {
	return r.ShowFunc(tournamentID, id)
}

func (r *mockRepository) Update(d Discipline) (Discipline, error) {
	return r.UpdateFunc(d)
}

func (r *mockRepository) Delete(tournamentID, id string) error {
	return r.DeleteFunc(tournamentID, id)
}

func newRepository() *mockRepository {
	return &mockRepository{
		AddFunc:    func(d Discipline) (Discipline, error) { return Discipline{}, nil },
		GetAllFunc: func(tournamentID string) ([]Discipline, error) { return []Discipline{}, nil },
		ShowFunc:   func(tournamentID, id string) (Discipline, error) { return Discipline{}, nil },
		UpdateFunc: func(d Discipline) (Discipline, error) { return Discipline{}, nil },
		DeleteFunc: func(tournamentID, id string) error { return nil },
	}
}

func TestServiceAdd_Success(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	called := false
	repo.AddFunc = func(d Discipline) (Discipline, error) {
		called = true
		return Discipline{}, nil
	}

	_, err := service.Add(Discipline{})

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !called {
		t.Fatal("Expected Add func to be called")
	}
}

func TestServiceAdd_Failure(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	repo.AddFunc = func(d Discipline) (Discipline, error) {
		return Discipline{}, errors.New("Error")
	}

	_, err := service.Add(Discipline{})

	if err == nil {
		t.Fatalf("Expected failure, got success")
	}
}

func TestServiceGetAll_Success(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	called := false
	repo.GetAllFunc = func(tournamentID string) ([]Discipline, error) {
		called = true
		return []Discipline{}, nil
	}

	_, err := service.GetAll("abc")

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !called {
		t.Fatal("Expected GetAll func to be called")
	}
}

func TestServiceGetAll_Failure(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	repo.GetAllFunc = func(tournamentID string) ([]Discipline, error) {
		return []Discipline{}, errors.New("Error")
	}

	_, err := service.GetAll("abc")

	if err == nil {
		t.Fatalf("Expected failire, got success")
	}
}

func TestServiceShow_Success(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	called := false
	repo.ShowFunc = func(tournamentID, id string) (Discipline, error) {
		called = true
		return Discipline{}, nil
	}

	_, err := service.Show("abc", "xyz")

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !called {
		t.Fatal("Expected Show func to be called")
	}
}

func TestServiceShow_Failure(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	repo.ShowFunc = func(tournamentID, id string) (Discipline, error) {
		return Discipline{}, errors.New("Error")
	}

	_, err := service.Show("abc", "xyz")

	if err == nil {
		t.Fatalf("Expected failure, got success")
	}
}

func TestServiceUpdate_Success(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	called := false
	repo.UpdateFunc = func(d Discipline) (Discipline, error) {
		called = true
		return Discipline{}, nil
	}

	_, err := service.Update(Discipline{})

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !called {
		t.Fatal("Expected Update func to be called")
	}
}

func TestServiceUpdate_Failure(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	repo.UpdateFunc = func(d Discipline) (Discipline, error) {
		return Discipline{}, errors.New("Error")
	}

	_, err := service.Update(Discipline{})

	if err == nil {
		t.Fatalf("Expected failure, got success")
	}
}

func TestServiceDelete_Success(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	called := false
	repo.DeleteFunc = func(tournamentID, id string) error {
		called = true
		return nil
	}

	err := service.Delete("abc", "xyz")

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !called {
		t.Fatal("Expected Show func to be called")
	}
}

func TestServiceDelete_Failure(t *testing.T) {
	repo := newRepository()
	service := NewService(repo)

	repo.DeleteFunc = func(tournamentID, id string) error {
		return errors.New("Error")
	}

	err := service.Delete("abc", "xyz")

	if err == nil {
		t.Fatalf("Expected failure, got success")
	}
}
