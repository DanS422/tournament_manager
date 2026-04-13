package tournament

import (
	"errors"
)

type ServiceInterface interface {
	List() ([]Tournament, error)
	Create(string, string) (Tournament, error)
	Show(int) (Tournament, error)
	Update(int, string, string) error
	Delete(int) error
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, location string) (Tournament, error) {
	if name == "" {
		return Tournament{}, errors.New("Name is required")
	}
	t, err := s.repo.Add(Tournament{Name: name, Location: location})
	if err != nil {
		return Tournament{}, err
	}

	return t, nil
}

func (s *Service) List() ([]Tournament, error) {
	tournaments, err := s.repo.GetAll()

	if err != nil {
		return []Tournament{}, err
	}

	return tournaments, nil
}

func (s *Service) Show(id int) (Tournament, error) {
	return s.repo.Show(id)
}

func (s *Service) Update(id int, name, location string) error {
	if name == "" {
		return errors.New("Name is required")
	}

	return s.repo.Update(id, Tournament{Name: name, Location: location})
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
