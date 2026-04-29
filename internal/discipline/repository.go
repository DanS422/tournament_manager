package discipline

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type RepositoryInterface interface {
	Add(d Discipline) (Discipline, error)
	GetAll(tournamentID string) ([]Discipline, error)
	Show(tournamentID, id string) (Discipline, error)
	Update(d Discipline) (Discipline, error)
	Delete(tournamentID, id string) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(d Discipline) (Discipline, error) {
	d.ID = uuid.NewString()
	_, err := r.db.Exec(
		"INSERT INTO disciplines (id, name, no_of_team_players, tournament_id) VALUES (?, ?, ?, ?)",
		d.ID, d.Name, d.NoOfTeamPlayers, d.TournamentID,
	)

	if err != nil {
		return Discipline{}, err
	}

	return d, nil
}

func (r *Repository) GetAll(tournamentID string) ([]Discipline, error) {
	rows, err := r.db.Query("SELECT id, name, no_of_team_players, tournament_id FROM disciplines")

	if err != nil {
		return []Discipline{}, err
	}

	defer rows.Close()

	list := []Discipline{}

	for rows.Next() {
		var d Discipline

		if err := rows.Scan(&d.ID, &d.Name, &d.NoOfTeamPlayers, &d.TournamentID); err != nil {
			return []Discipline{}, err
		}

		list = append(list, d)
	}

	if err := rows.Err(); err != nil {
		return []Discipline{}, err
	}

	return list, nil

}

func (r *Repository) Show(tournamentID, id string) (Discipline, error) {
	var d Discipline

	err := r.db.QueryRow(
		"SELECT id, name, no_of_team_players, tournament_id FROM disciplines where tournament_id = ? AND id = ?",
		tournamentID, id,
	).Scan(&d.ID, &d.Name, &d.NoOfTeamPlayers, &d.TournamentID)

	if err != nil {
		return Discipline{}, err
	}

	return d, nil
}

func (r *Repository) Update(d Discipline) (Discipline, error) {
	result, err := r.db.Exec(
		"UPDATE disciplines SET name = ?, no_of_team_players = ? WHERE id = ? AND tournament_id = ?",
		d.Name, d.NoOfTeamPlayers, d.ID, d.TournamentID,
	)

	if err != nil {
		return Discipline{}, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return Discipline{}, err
	}

	if rowsAffected == 0 {
		return Discipline{}, errors.New("record not found")
	}

	return Discipline{}, nil
}

func (r *Repository) Delete(tournamentID, id string) error {
	result, err := r.db.Exec(
		"DELETE FROM disciplines where tournament_id = ? AND id = ?",
		tournamentID, id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}
