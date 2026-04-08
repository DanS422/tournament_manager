package tournament

import "errors"

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
	return s.repo.Add(Tournament{Name: name, Location: location}), nil
}

func (s *Service) List() []Tournament {
	return s.repo.GetAll()
}

func (s *Service) show(id int) (Tournament, error) {
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
