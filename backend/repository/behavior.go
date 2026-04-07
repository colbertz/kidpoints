package repository

import (
	"points-app/database"
	"points-app/models"
)

type BehaviorRepository struct{}

func NewBehaviorRepository() *BehaviorRepository {
	return &BehaviorRepository{}
}

func (r *BehaviorRepository) GetAll() ([]models.Behavior, error) {
	rows, err := database.DB.Query("SELECT id, name, type, points FROM behaviors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var behaviors []models.Behavior
	for rows.Next() {
		var b models.Behavior
		if err := rows.Scan(&b.ID, &b.Name, &b.Type, &b.Points); err != nil {
			return nil, err
		}
		behaviors = append(behaviors, b)
	}
	return behaviors, nil
}

func (r *BehaviorRepository) GetByID(id int) (*models.Behavior, error) {
	var b models.Behavior
	err := database.DB.QueryRow("SELECT id, name, type, points FROM behaviors WHERE id = ?", id).
		Scan(&b.ID, &b.Name, &b.Type, &b.Points)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BehaviorRepository) Create(b *models.Behavior) error {
	result, err := database.DB.Exec(
		"INSERT INTO behaviors (name, type, points) VALUES (?, ?, ?)",
		b.Name, b.Type, b.Points,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	b.ID = int(id)
	return nil
}

func (r *BehaviorRepository) Update(b *models.Behavior) error {
	_, err := database.DB.Exec(
		"UPDATE behaviors SET name = ?, type = ?, points = ? WHERE id = ?",
		b.Name, b.Type, b.Points, b.ID,
	)
	return err
}

func (r *BehaviorRepository) Delete(id int) error {
	_, err := database.DB.Exec("DELETE FROM behaviors WHERE id = ?", id)
	return err
}
