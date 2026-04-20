package player

type Player struct {
	ID           string `json:"id"`
	Age          int
	Gender       string `json:"gender" validate:"oneof=male female diverse"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	TournamentID string `json:"tournament_id" validate:"required"`
}
