package discipline

type Discipline struct {
	ID              string `json:"id"`
	Name            string `json:"name" validate:"required"`
	NoOfTeamPlayers int    `json:"no_of_team_players" validate:"required,min=1"`
	TournamentID    string `json:"tournament_id" validate:"required"`
}
