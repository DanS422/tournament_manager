package player

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type RepositoryInterface interface {
	Add(p Player) (Player, error)
	GetAll(tournamentID string) ([]Player, error)
	Show(id string) (Player, error)
	Update(p Player) error
	Delete(id string) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(p Player) (Player, error) {
	p.ID = uuid.NewString()

	_, err := r.db.Exec(
		"INSERT INTO players (id, first_name, last_name, gender, tournament_id) VALUES (?, ?, ?, ?, ?)",
		p.ID, p.FirstName, p.LastName, p.Gender, p.TournamentID,
	)

	if err != nil {
		return Player{}, err
	}

	return p, nil
}

func (r *Repository) GetAll(tournamentID string) ([]Player, error) {
	rows, err := r.db.Query(
		"SELECT id, first_name, last_name, gender, tournament_id FROM players WHERE tournament_id = ?",
		tournamentID,
	)

	if err != nil {
		return []Player{}, err
	}

	defer rows.Close()

	list := []Player{}
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Gender, &p.TournamentID); err != nil {
			return []Player{}, err
		}

		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return []Player{}, err
	}

	return list, nil
}

func (r *Repository) Show(id string) (Player, error) {
	var p Player

	err := r.db.QueryRow(
		"SELECT id, first_name, last_name, gender, tournament_id FROM players WHERE id = ?",
		id,
	).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Gender, &p.TournamentID)

	if err == sql.ErrNoRows {
		return Player{}, errors.New("player not found")
	}

	if err != nil {
		return Player{}, err
	}

	return p, nil
}

func (r *Repository) Update(p Player) error {
	result, err := r.db.Exec(
		"UPDATE players SET first_name = ?, last_name = ?, gender = ? WHERE id = ?",
		p.FirstName, p.LastName, p.Gender, p.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("player not found")
	}

	return nil
}

func (r *Repository) Delete(id string) error {
	result, err := r.db.Exec(
		"DELETE FROM players WHERE id = ?",
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("player not found")
	}

	return nil
}
