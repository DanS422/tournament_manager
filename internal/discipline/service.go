package discipline

type ServiceInterface interface {
	Add(d Discipline) (Discipline, error)
	GetAll(tournamentID string) ([]Discipline, error)
	Show(tournamentID, id string) (Discipline, error)
	Update(d Discipline) (Discipline, error)
	Delete(tournamentID, id string) error
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(d Discipline) (Discipline, error) {
	return s.repo.Add(d)
}

func (s *Service) GetAll(tournamentID string) ([]Discipline, error) {
	return s.repo.GetAll(tournamentID)
}

func (s *Service) Show(tournamentID, id string) (Discipline, error) {
	return s.repo.Show(tournamentID, id)
}

func (s *Service) Update(d Discipline) (Discipline, error) {
	return s.repo.Update(d)
}

func (s *Service) Delete(tournamentID, id string) error {
	return s.repo.Delete(tournamentID, id)
}
