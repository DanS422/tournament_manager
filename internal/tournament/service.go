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
	return s.repo.Add(t)
}

func (s *Service) List() ([]Tournament, error) {
	return s.repo.GetAll()
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
