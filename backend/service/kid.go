package service

import (
	"errors"

	"points-app/models"
	"points-app/repository"
)

var ErrKidNotFound = errors.New("kid not found")
var ErrInsufficientPoints = errors.New("insufficient points")

type KidService struct {
	kidRepo     *repository.KidRepository
	behaviorRepo *repository.BehaviorRepository
	recordRepo   *repository.RecordRepository
}

func NewKidService(kidRepo *repository.KidRepository, behaviorRepo *repository.BehaviorRepository, recordRepo *repository.RecordRepository) *KidService {
	return &KidService{
		kidRepo:     kidRepo,
		behaviorRepo: behaviorRepo,
		recordRepo:   recordRepo,
	}
}

func (s *KidService) GetAllKids() ([]models.Kid, error) {
	return s.kidRepo.GetAll()
}

func (s *KidService) GetKidPoints(kidID int) (int, error) {
	kid, err := s.kidRepo.GetByID(kidID)
	if err != nil {
		return 0, ErrKidNotFound
	}
	return kid.Points, nil
}

func (s *KidService) GetKidByName(name string) (*models.Kid, error) {
	return s.kidRepo.GetByName(name)
}

func (s *KidService) GetKid(kidID int) (*models.Kid, error) {
	return s.kidRepo.GetByID(kidID)
}

func (s *KidService) AddPoints(kidID, behaviorID int) error {
	behavior, err := s.behaviorRepo.GetByID(behaviorID)
	if err != nil {
		return errors.New("behavior not found")
	}

	if behavior.Type != models.BehaviorTypeAdd {
		return errors.New("behavior is not an add type")
	}

	if err := s.kidRepo.UpdatePoints(kidID, behavior.Points); err != nil {
		return err
	}

	return s.recordRepo.Create(kidID, behaviorID, behavior.Points)
}

func (s *KidService) SubPoints(kidID, behaviorID int) error {
	behavior, err := s.behaviorRepo.GetByID(behaviorID)
	if err != nil {
		return errors.New("behavior not found")
	}

	if behavior.Type != models.BehaviorTypeSub {
		return errors.New("behavior is not a sub type")
	}

	kid, err := s.kidRepo.GetByID(kidID)
	if err != nil {
		return ErrKidNotFound
	}

	if kid.Points < behavior.Points {
		return ErrInsufficientPoints
	}

	if err := s.kidRepo.UpdatePoints(kidID, -behavior.Points); err != nil {
		return err
	}

	return s.recordRepo.Create(kidID, behaviorID, -behavior.Points)
}

func (s *KidService) CreateKid(kid *models.Kid) error {
	if kid.Points < 0 {
		kid.Points = 0
	}
	return s.kidRepo.Create(kid)
}

func (s *KidService) UpdateKid(kid *models.Kid) error {
	_, err := s.kidRepo.GetByID(kid.ID)
	if err != nil {
		return ErrKidNotFound
	}
	return s.kidRepo.Update(kid)
}

func (s *KidService) DeleteKid(id int) error {
	_, err := s.kidRepo.GetByID(id)
	if err != nil {
		return ErrKidNotFound
	}
	return s.kidRepo.Delete(id)
}
