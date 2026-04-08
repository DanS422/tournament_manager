package tournament

import "errors"

type Repository struct {
	data   map[int]Tournament
	nextID int
}

func NewRepository() *Repository {
	return &Repository{
		data:   make(map[int]Tournament),
		nextID: 1,
	}
}

func (r *Repository) Add(t Tournament) Tournament {
	t.ID = r.nextID
	r.nextID++
	r.data[t.ID] = t
	return t
}

func (r *Repository) GetAll() []Tournament {
	list := []Tournament{}
	for _, t := range r.data {
		list = append(list, t)
	}

	return list
}

func (r *Repository) Show(id int) (Tournament, error) {
	tournament, ok := r.data[id]

	if !ok {
		return Tournament{}, errors.New("tournament not found")
	}

	return tournament, nil
}

func (r *Repository) Update(id int, t Tournament) error {
	_, ok := r.data[id]

	if !ok {
		return errors.New("tournament not found")
	}

	t.ID = id
	r.data[id] = t
	return nil
}

func (r *Repository) Delete(id int) error {
	_, ok := r.data[id]

	if !ok {
		return errors.New("tournament not found")
	}

	delete(r.data, id)
	return nil
}
