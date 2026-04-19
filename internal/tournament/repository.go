package tournament

import (
	"database/sql"
	"errors"
)

type RepositoryInterface interface {
	Add(t Tournament) (Tournament, error)
	GetAll() ([]Tournament, error)
	Show(id int) (Tournament, error)
	Update(id int, t Tournament) error
	Delete(id int) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(t Tournament) (Tournament, error) {
	result, err := r.db.Exec(
		"INSERT INTO tournaments (name, location) VALUES (?, ?)",
		t.Name, t.Location,
	)

	if err != nil {
		return Tournament{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return Tournament{}, err
	}

	t.ID = int(id)
	return t, nil
}

func (r *Repository) GetAll() ([]Tournament, error) {
	rows, err := r.db.Query("SELECT id, name, location FROM tournaments")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := []Tournament{}
	for rows.Next() {
		var t Tournament
		if err := rows.Scan(&t.ID, &t.Name, &t.Location); err != nil {
			return nil, err
		}
		list = append(list, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Repository) Show(id int) (Tournament, error) {
	var t Tournament

	err := r.db.QueryRow("SELECT id, name, location FROM tournaments WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Location)
	if err == sql.ErrNoRows {
		return Tournament{}, errors.New("tournament not found")
	}

	if err != nil {
		return Tournament{}, err
	}

	return t, nil
}

func (r *Repository) Update(id int, t Tournament) error {
	result, err := r.db.Exec(
		"UPDATE tournaments SET name = ?, location = ? where id = ?",
		t.Name,
		t.Location,
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
		return errors.New("tournament not found")
	}

	return nil
}

func (r *Repository) Delete(id int) error {
	result, err := r.db.Exec(
		"DELETE FROM tournaments WHERE id = ?",
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
		return errors.New("tournament not found")
	}

	return nil
}
