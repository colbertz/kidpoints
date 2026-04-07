package service

import (
	"testing"

	"points-app/database"
	"points-app/models"
	"points-app/repository"
)

func setupTestDB(t *testing.T) {
	tmpFile := t.TempDir() + "/test.db"
	database.Init(tmpFile)
	t.Cleanup(func() { database.Close() })
}

func TestKidService_GetAllKids(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	kids, err := svc.GetAllKids()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kids) != 2 {
		t.Errorf("expected 2 kids, got %d", len(kids))
	}
}

func TestKidService_AddPoints(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	err := svc.AddPoints(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	points, err := svc.GetKidPoints(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if points != 2 {
		t.Errorf("expected 2 points, got %d", points)
	}
}

func TestKidService_SubPoints(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	kidRepo.UpdatePoints(1, 10)

	err := svc.SubPoints(1, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	points, err := svc.GetKidPoints(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if points != 8 {
		t.Errorf("expected 8 points, got %d", points)
	}
}

func TestKidService_SubPoints_InsufficientPoints(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	err := svc.SubPoints(1, 5)
	if err != ErrInsufficientPoints {
		t.Errorf("expected ErrInsufficientPoints, got %v", err)
	}
}

func TestKidService_AddPoints_WrongBehaviorType(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	err := svc.AddPoints(1, 5)
	if err == nil {
		t.Error("expected error for wrong behavior type")
	}
}

func TestKidService_CreateKid(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	kid := &models.Kid{Name: "TestKid", Avatar: "🧒", Points: 100}
	err := svc.CreateKid(kid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kid.ID == 0 {
		t.Error("expected kid ID to be set")
	}

	kids, _ := svc.GetAllKids()
	if len(kids) != 3 {
		t.Errorf("expected 3 kids, got %d", len(kids))
	}
}

func TestKidService_CreateKid_NegativePoints(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	kid := &models.Kid{Name: "TestKid", Avatar: "🧒", Points: -50}
	err := svc.CreateKid(kid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kid.Points != 0 {
		t.Errorf("expected 0 points, got %d", kid.Points)
	}
}

func TestKidService_UpdateKid(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	kid := &models.Kid{ID: 1, Name: "UpdatedKid", Avatar: "👦", Points: 200}
	err := svc.UpdateKid(kid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	kidBack, _ := kidRepo.GetByID(1)
	if kidBack.Name != "UpdatedKid" {
		t.Errorf("expected name 'UpdatedKid', got '%s'", kidBack.Name)
	}
	if kidBack.Points != 200 {
		t.Errorf("expected 200 points, got %d", kidBack.Points)
	}
}

func TestKidService_UpdateKid_NotFound(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	kid := &models.Kid{ID: 999, Name: "NotExist", Avatar: "👻", Points: 0}
	err := svc.UpdateKid(kid)
	if err != ErrKidNotFound {
		t.Errorf("expected ErrKidNotFound, got %v", err)
	}
}

func TestKidService_DeleteKid(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	err := svc.DeleteKid(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	kids, _ := svc.GetAllKids()
	if len(kids) != 1 {
		t.Errorf("expected 1 kid, got %d", len(kids))
	}
}

func TestKidService_DeleteKid_NotFound(t *testing.T) {
	setupTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	svc := NewKidService(kidRepo, behaviorRepo, recordRepo)

	err := svc.DeleteKid(999)
	if err != ErrKidNotFound {
		t.Errorf("expected ErrKidNotFound, got %v", err)
	}
}
