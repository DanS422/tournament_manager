package attendant

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type RepositoryInterface interface {
	Add(a Attendant) (Attendant, error)
	GetAll(tournamentID string) ([]DisplayAttendant, error)
	Show(tournamentID, id string) (DisplayAttendant, error)
	Delete(tournamentID, id string) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(a Attendant) (Attendant, error) {
	a.ID = uuid.NewString()

	_, err := r.db.Exec(
		"INSERT INTO attendants (id, player_id, tournament_id) VALUES (?, ?, ?)",
		a.ID, a.PlayerID, a.TournamentID,
	)

	if err != nil {
		return Attendant{}, err
	}

	return a, err
}

func (r *Repository) GetAll(tournamentID string) ([]DisplayAttendant, error) {
	rows, err := r.db.Query(
		`SELECT
			a.id,
			a.player_id,
			a.tournament_id,
			p.first_name,
			p.last_name,
			p.gender
		FROM attendants a
		JOIN players p ON p.id = a.player_id
		WHERE a.tournament_id = ?
		ORDER BY p.last_name, p.first_name`,
		tournamentID,
	)
	if err != nil {
		return []DisplayAttendant{}, err
	}
	defer rows.Close()

	list := []DisplayAttendant{}
	for rows.Next() {
		var a DisplayAttendant
		if err := rows.Scan(
			&a.ID,
			&a.PlayerID,
			&a.TournamentID,
			&a.FirstName,
			&a.LastName,
			&a.Gender,
		); err != nil {
			return []DisplayAttendant{}, err
		}

		list = append(list, a)
	}

	if err := rows.Err(); err != nil {
		return []DisplayAttendant{}, err
	}

	return list, nil
}

func (r *Repository) Show(tournamentID, id string) (DisplayAttendant, error) {
	var a DisplayAttendant

	err := r.db.QueryRow(
		`SELECT
			a.id,
			a.player_id,
			a.tournament_id,
			p.first_name,
			p.last_name,
			p.gender
		FROM attendants a
		JOIN players p on p.id = a.player_id
		WHERE a.tournament_id = ?`,
		tournamentID,
	).Scan(
		&a.ID,
		&a.PlayerID,
		&a.TournamentID,
		&a.FirstName,
		&a.LastName,
		&a.Gender,
	)

	if err != nil {
		return DisplayAttendant{}, err
	}

	return a, nil
}

func (r *Repository) Delete(tournamentID, id string) error {
	result, err := r.db.Exec(
		"DELETE FROM attendants WHERE id = ? AND tournament_id = ?",
		id, tournamentID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("attendant not found")
	}

	return nil
}
