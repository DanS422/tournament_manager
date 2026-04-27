package attendant

type Attendant struct {
	ID           string `json:"id"`
	PlayerID     string `json:"player_id" validate:"required"`
	TournamentID string `json:"tournament_id" validate:"required"`
}

type DisplayAttendant struct {
	ID           string `json:"id"`
	PlayerID     string `json:"player_id"`
	TournamentID string `json:"tournament_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
}

func (r *DisplayAttendant) FullName() string {
	return r.LastName + ", " + r.FirstName
}
