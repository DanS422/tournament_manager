package attendant

type ServiceInterface interface {
	Add(a Attendant) (Attendant, error)
	GetAll(tournamentID string) ([]DisplayAttendant, error)
	Show(tournamentID, id string) (DisplayAttendant, error)
	Delete(tournamentID, id string) error
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(a Attendant) (Attendant, error) {
	return s.repo.Add(a)
}

func (s *Service) GetAll(tournamentID string) ([]DisplayAttendant, error) {
	return s.repo.GetAll(tournamentID)
}

func (s *Service) Show(tournamentID, id string) (DisplayAttendant, error) {
	return s.repo.Show(tournamentID, id)
}

func (s *Service) Delete(tournamentID, id string) error {
	return s.repo.Delete(tournamentID, id)
}
