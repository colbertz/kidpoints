package models

type Kid struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Points  int    `json:"points"`
}
