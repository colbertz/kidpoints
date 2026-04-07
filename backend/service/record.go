package service

import (
	"database/sql"
	"errors"

	"points-app/models"
	"points-app/repository"
)

type RecordService struct {
	recordRepo *repository.RecordRepository
	kidRepo    *repository.KidRepository
}

func NewRecordService(recordRepo *repository.RecordRepository, kidRepo *repository.KidRepository) *RecordService {
	return &RecordService{
		recordRepo: recordRepo,
		kidRepo:    kidRepo,
	}
}

func (s *RecordService) GetRecordsByKid(kidID int) ([]models.Record, error) {
	kid, err := s.kidRepo.GetByID(kidID)
	if err != nil {
		return nil, ErrKidNotFound
	}
	_ = kid

	return s.recordRepo.GetByKidID(kidID)
}

func (s *RecordService) GetAllRecords() ([]models.Record, error) {
	return s.recordRepo.GetAll()
}

var ErrRecordNotFound = errors.New("record not found")
var ErrRecordAlreadyReversed = errors.New("record already reversed")

func (s *RecordService) ReverseRecord(recordID int, reason string) (*models.Record, error) {
	record, err := s.recordRepo.GetByID(recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	if record.Reversed {
		return nil, ErrRecordAlreadyReversed
	}

	if err := s.kidRepo.UpdatePoints(record.KidID, -record.Points); err != nil {
		return nil, err
	}

	if err := s.recordRepo.Reverse(recordID, reason); err != nil {
		return nil, err
	}

	return s.recordRepo.GetByID(recordID)
}
