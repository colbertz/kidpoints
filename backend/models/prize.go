package models

type Prize struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Probability float64 `json:"probability"`
}
