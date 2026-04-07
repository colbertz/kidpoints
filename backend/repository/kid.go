package repository

import (
	"database/sql"

	"points-app/database"
	"points-app/models"
)

type KidRepository struct{}

func NewKidRepository() *KidRepository {
	return &KidRepository{}
}

func (r *KidRepository) GetAll() ([]models.Kid, error) {
	rows, err := database.DB.Query("SELECT id, name, avatar, points FROM kids")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kids []models.Kid
	for rows.Next() {
		var k models.Kid
		if err := rows.Scan(&k.ID, &k.Name, &k.Avatar, &k.Points); err != nil {
			return nil, err
		}
		kids = append(kids, k)
	}
	return kids, nil
}

func (r *KidRepository) GetByID(id int) (*models.Kid, error) {
	var k models.Kid
	err := database.DB.QueryRow("SELECT id, name, avatar, points FROM kids WHERE id = ?", id).
		Scan(&k.ID, &k.Name, &k.Avatar, &k.Points)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *KidRepository) GetByName(name string) (*models.Kid, error) {
	var k models.Kid
	err := database.DB.QueryRow("SELECT id, name, avatar, points FROM kids WHERE name = ?", name).
		Scan(&k.ID, &k.Name, &k.Avatar, &k.Points)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *KidRepository) UpdatePoints(id, points int) error {
	_, err := database.DB.Exec("UPDATE kids SET points = points + ? WHERE id = ?", points, id)
	return err
}

func (r *KidRepository) GetPoints(id int) (int, error) {
	var points int
	err := database.DB.QueryRow("SELECT points FROM kids WHERE id = ?", id).Scan(&points)
	return points, err
}

func (r *KidRepository) DeductPoints(id, points int) error {
	result, err := database.DB.Exec("UPDATE kids SET points = points - ? WHERE id = ? AND points >= ?", points, id, points)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *KidRepository) Create(kid *models.Kid) error {
	result, err := database.DB.Exec(
		"INSERT INTO kids (name, avatar, points) VALUES (?, ?, ?)",
		kid.Name, kid.Avatar, kid.Points,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	kid.ID = int(id)
	return nil
}

func (r *KidRepository) Update(kid *models.Kid) error {
	_, err := database.DB.Exec(
		"UPDATE kids SET name = ?, avatar = ?, points = ? WHERE id = ?",
		kid.Name, kid.Avatar, kid.Points, kid.ID,
	)
	return err
}

func (r *KidRepository) Delete(id int) error {
	_, err := database.DB.Exec("DELETE FROM kids WHERE id = ?", id)
	return err
}
