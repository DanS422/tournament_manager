package tournament

type Tournament struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}
