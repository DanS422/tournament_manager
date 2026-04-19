package tournament

type ServiceInterface interface {
	List() ([]Tournament, error)
	Create(t Tournament) (Tournament, error)
	Show(string) (Tournament, error)
	Update(t Tournament) error
	Delete(string) error
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(t Tournament) (Tournament, error) {
	t, err := s.repo.Add(t)
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

func (s *Service) Show(id string) (Tournament, error) {
	return s.repo.Show(id)
}

func (s *Service) Update(t Tournament) error {
	return s.repo.Update(t)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
