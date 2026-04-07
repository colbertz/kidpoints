package service

import (
	"errors"
	"math/rand"

	"points-app/models"
	"points-app/repository"
)

var ErrNoPrizesAvailable = errors.New("no prizes available")
var ErrInsufficientPointsForDraw = errors.New("insufficient points for draw")

const PointsPerDraw = 2

type DrawService struct {
	kidRepo   *repository.KidRepository
	prizeRepo *repository.PrizeRepository
}

func NewDrawService(kidRepo *repository.KidRepository, prizeRepo *repository.PrizeRepository) *DrawService {
	return &DrawService{
		kidRepo:   kidRepo,
		prizeRepo: prizeRepo,
	}
}

func (s *DrawService) Draw(kidID int) (*models.Prize, error) {
	kid, err := s.kidRepo.GetByID(kidID)
	if err != nil {
		return nil, ErrKidNotFound
	}

	if kid.Points < PointsPerDraw {
		return nil, ErrInsufficientPointsForDraw
	}

	prizes, err := s.prizeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(prizes) == 0 {
		return nil, ErrNoPrizesAvailable
	}

	winner := selectWinner(prizes)

	s.kidRepo.DeductPoints(kidID, PointsPerDraw)

	return winner, nil
}

func selectWinner(prizes []models.Prize) *models.Prize {
	totalProbability := 0.0
	for _, p := range prizes {
		totalProbability += p.Probability / 100.0
	}

	r := rand.Float64() * totalProbability

	cumulative := 0.0
	for i := range prizes {
		cumulative += prizes[i].Probability / 100.0
		if r <= cumulative {
			return &prizes[i]
		}
	}

	return &prizes[len(prizes)-1]
}
