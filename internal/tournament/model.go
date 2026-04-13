package tournament

type Tournament struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}
