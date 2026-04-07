package service

import (
	"testing"

	"points-app/database"
	"points-app/models"
	"points-app/repository"
)

func setupDrawTestDB(t *testing.T) {
	tmpFile := t.TempDir() + "/test.db"
	database.Init(tmpFile)
	t.Cleanup(func() { database.Close() })

	prizes := []models.Prize{
		{ID: 1, Name: "一等奖", Icon: "🎉", Probability: 10},
		{ID: 2, Name: "二等奖", Icon: "🎁", Probability: 20},
		{ID: 3, Name: "三等奖", Icon: "🎀", Probability: 30},
		{ID: 4, Name: "再接再厉", Icon: "🍀", Probability: 40},
	}
	for _, p := range prizes {
		repository.NewPrizeRepository().Create(&p)
	}
}

func TestDrawService_InsufficientPoints(t *testing.T) {
	setupDrawTestDB(t)
	kidRepo := repository.NewKidRepository()
	prizeRepo := repository.NewPrizeRepository()
	svc := NewDrawService(kidRepo, prizeRepo)

	_, err := svc.Draw(1)
	if err != ErrInsufficientPointsForDraw {
		t.Errorf("expected ErrInsufficientPointsForDraw, got %v", err)
	}
}

func TestDrawService_Success(t *testing.T) {
	setupDrawTestDB(t)
	kidRepo := repository.NewKidRepository()
	prizeRepo := repository.NewPrizeRepository()
	svc := NewDrawService(kidRepo, prizeRepo)

	kidRepo.UpdatePoints(1, 10)

	prize, err := svc.Draw(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prize == nil {
		t.Fatal("expected prize, got nil")
	}

	points, _ := kidRepo.GetPoints(1)
	if points != 8 {
		t.Errorf("expected 8 points after draw, got %d", points)
	}
}

func TestDrawService_KidNotFound(t *testing.T) {
	setupDrawTestDB(t)
	kidRepo := repository.NewKidRepository()
	prizeRepo := repository.NewPrizeRepository()
	svc := NewDrawService(kidRepo, prizeRepo)

	_, err := svc.Draw(999)
	if err != ErrKidNotFound {
		t.Errorf("expected ErrKidNotFound, got %v", err)
	}
}

func TestSelectWinner(t *testing.T) {
	prizes := []models.Prize{
		{Name: "p1", Probability: 50},
		{Name: "p2", Probability: 50},
	}

	for i := 0; i < 100; i++ {
		winner := selectWinner(prizes)
		if winner == nil {
			t.Fatal("expected winner")
		}
	}
}
