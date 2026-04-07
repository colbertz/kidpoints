package service

import (
	"testing"

	"points-app/database"
	"points-app/repository"
)

func setupRecordTestDB(t *testing.T) {
	tmpFile := t.TempDir() + "/test.db"
	database.Init(tmpFile)
	t.Cleanup(func() { database.Close() })
}

func TestRecordService_GetAllRecords(t *testing.T) {
	setupRecordTestDB(t)
	recordRepo := repository.NewRecordRepository()
	kidRepo := repository.NewKidRepository()
	svc := NewRecordService(recordRepo, kidRepo)

	records, err := svc.GetAllRecords()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

func TestRecordService_GetRecordsByKid(t *testing.T) {
	setupRecordTestDB(t)
	recordRepo := repository.NewRecordRepository()
	kidRepo := repository.NewKidRepository()
	kidSvc := NewKidService(kidRepo, repository.NewBehaviorRepository(), recordRepo)
	recordSvc := NewRecordService(recordRepo, kidRepo)

	kidSvc.AddPoints(1, 1)

	records, err := recordSvc.GetRecordsByKid(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}
}
