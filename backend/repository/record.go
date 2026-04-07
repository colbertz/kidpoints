package repository

import (
	"time"

	"points-app/database"
	"points-app/models"
)

type RecordRepository struct{}

func NewRecordRepository() *RecordRepository {
	return &RecordRepository{}
}

func (r *RecordRepository) Create(kidID, behaviorID, points int) error {
	_, err := database.DB.Exec(
		"INSERT INTO records (kid_id, behavior_id, points) VALUES (?, ?, ?)",
		kidID, behaviorID, points,
	)
	return err
}

func (r *RecordRepository) GetByID(id int) (*models.Record, error) {
	var rec models.Record
	err := database.DB.QueryRow(`
		SELECT r.id, r.kid_id, r.behavior_id, b.name, r.points, r.created_at, r.reversed, COALESCE(r.reversed_reason, ''), r.reversed_at
		FROM records r
		JOIN behaviors b ON r.behavior_id = b.id
		WHERE r.id = ?
	`, id).Scan(&rec.ID, &rec.KidID, &rec.BehaviorID, &rec.BehaviorName, &rec.Points, &rec.CreatedAt, &rec.Reversed, &rec.ReversedReason, &rec.ReversedAt)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecordRepository) GetByKidID(kidID int) ([]models.Record, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.kid_id, r.behavior_id, b.name, r.points, r.created_at, r.reversed, COALESCE(r.reversed_reason, ''), r.reversed_at
		FROM records r
		JOIN behaviors b ON r.behavior_id = b.id
		WHERE r.kid_id = ?
		ORDER BY r.created_at DESC
	`, kidID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var rec models.Record
		if err := rows.Scan(&rec.ID, &rec.KidID, &rec.BehaviorID, &rec.BehaviorName, &rec.Points, &rec.CreatedAt, &rec.Reversed, &rec.ReversedReason, &rec.ReversedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *RecordRepository) GetAll() ([]models.Record, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.kid_id, r.behavior_id, b.name, r.points, r.created_at, r.reversed, COALESCE(r.reversed_reason, ''), r.reversed_at
		FROM records r
		JOIN behaviors b ON r.behavior_id = b.id
		ORDER BY r.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var rec models.Record
		if err := rows.Scan(&rec.ID, &rec.KidID, &rec.BehaviorID, &rec.BehaviorName, &rec.Points, &rec.CreatedAt, &rec.Reversed, &rec.ReversedReason, &rec.ReversedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *RecordRepository) Reverse(id int, reason string) error {
	now := time.Now()
	_, err := database.DB.Exec(
		"UPDATE records SET reversed = TRUE, reversed_reason = ?, reversed_at = ? WHERE id = ?",
		reason, now, id,
	)
	return err
}
