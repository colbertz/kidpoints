package repository

import (
	"points-app/database"
	"points-app/models"
)

type PrizeRepository struct{}

func NewPrizeRepository() *PrizeRepository {
	return &PrizeRepository{}
}

func (r *PrizeRepository) GetAll() ([]models.Prize, error) {
	rows, err := database.DB.Query("SELECT id, name, icon, probability FROM prizes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prizes []models.Prize
	for rows.Next() {
		var p models.Prize
		if err := rows.Scan(&p.ID, &p.Name, &p.Icon, &p.Probability); err != nil {
			return nil, err
		}
		prizes = append(prizes, p)
	}
	return prizes, nil
}

func (r *PrizeRepository) GetByID(id int) (*models.Prize, error) {
	var p models.Prize
	err := database.DB.QueryRow("SELECT id, name, icon, probability FROM prizes WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Icon, &p.Probability)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PrizeRepository) Create(p *models.Prize) error {
	result, err := database.DB.Exec(
		"INSERT INTO prizes (name, icon, probability) VALUES (?, ?, ?)",
		p.Name, p.Icon, p.Probability,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)
	return nil
}

func (r *PrizeRepository) Update(p *models.Prize) error {
	_, err := database.DB.Exec(
		"UPDATE prizes SET name = ?, icon = ?, probability = ? WHERE id = ?",
		p.Name, p.Icon, p.Probability, p.ID,
	)
	return err
}

func (r *PrizeRepository) Delete(id int) error {
	_, err := database.DB.Exec("DELETE FROM prizes WHERE id = ?", id)
	return err
}
