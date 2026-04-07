package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"points-app/database"
	"points-app/models"
	"points-app/repository"
	"points-app/service"
)

func setupHandlerTestDB(t *testing.T) {
	tmpFile := t.TempDir() + "/test.db"
	database.Init(tmpFile)
	t.Cleanup(func() { database.Close() })
}

func TestKidHandler_GetKids(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.GET("/kids", handler.GetKids)

	req, _ := http.NewRequest("GET", "/kids", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var kids []interface{}
	json.Unmarshal(w.Body.Bytes(), &kids)
	if len(kids) != 2 {
		t.Errorf("expected 2 kids, got %d", len(kids))
	}
}

func TestKidHandler_AddPoints(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.POST("/kids/:id/points/add", handler.AddPoints)

	body, _ := json.Marshal(map[string]int{"id": 1, "behavior_id": 1})
	req, _ := http.NewRequest("POST", "/kids/1/points/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestKidHandler_SubPoints_InsufficientPoints(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.POST("/kids/:id/points/sub", handler.SubPoints)

	body, _ := json.Marshal(map[string]int{"id": 1, "behavior_id": 5})
	req, _ := http.NewRequest("POST", "/kids/1/points/sub", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestKidHandler_GetPoints(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.GET("/kids/:id/points", handler.GetPoints)

	req, _ := http.NewRequest("GET", "/kids/1/points", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestKidHandler_CreateKid(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.POST("/kids", handler.CreateKid)

	body, _ := json.Marshal(map[string]interface{}{"name": "NewKid", "avatar": "🧒", "points": 50})
	req, _ := http.NewRequest("POST", "/kids", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var kid models.Kid
	json.Unmarshal(w.Body.Bytes(), &kid)
	if kid.Name != "NewKid" {
		t.Errorf("expected name 'NewKid', got '%s'", kid.Name)
	}
}

func TestKidHandler_CreateKid_WithoutAvatar(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.POST("/kids", handler.CreateKid)

	body, _ := json.Marshal(map[string]interface{}{"name": "NoAvatarKid", "points": 0})
	req, _ := http.NewRequest("POST", "/kids", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var kid models.Kid
	json.Unmarshal(w.Body.Bytes(), &kid)
	if kid.Avatar != "👤" {
		t.Errorf("expected default avatar '👤', got '%s'", kid.Avatar)
	}
}

func TestKidHandler_CreateKid_MissingName(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.POST("/kids", handler.CreateKid)

	body, _ := json.Marshal(map[string]interface{}{"points": 0})
	req, _ := http.NewRequest("POST", "/kids", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestKidHandler_UpdateKid(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.PUT("/kids/:id", handler.UpdateKid)

	body, _ := json.Marshal(map[string]interface{}{"name": "UpdatedKid", "avatar": "👦", "points": 100})
	req, _ := http.NewRequest("PUT", "/kids/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestKidHandler_UpdateKid_NotFound(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.PUT("/kids/:id", handler.UpdateKid)

	body, _ := json.Marshal(map[string]interface{}{"name": "Ghost", "avatar": "👻", "points": 0})
	req, _ := http.NewRequest("PUT", "/kids/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestKidHandler_DeleteKid(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.DELETE("/kids/:id", handler.DeleteKid)

	req, _ := http.NewRequest("DELETE", "/kids/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestKidHandler_DeleteKid_NotFound(t *testing.T) {
	setupHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	recordRepo := repository.NewRecordRepository()
	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	handler := NewKidHandler(kidService)

	r := gin.New()
	r.DELETE("/kids/:id", handler.DeleteKid)

	req, _ := http.NewRequest("DELETE", "/kids/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
