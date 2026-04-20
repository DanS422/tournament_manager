package player

type ServiceInterface interface {
	Create(p Player) (Player, error)
	List() ([]Player, error)
	Show(id string) (Player, error)
	Update(p Player) error
	Delete(id string) error
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(p Player) (Player, error) {
	return s.repo.Add(p)
}

func (s *Service) List() ([]Player, error) {
	return s.repo.GetAll()
}

func (s *Service) Show(id string) (Player, error) {
	return s.repo.Show(id)
}

func (s *Service) Update(p Player) error {
	return s.repo.Update(p)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
