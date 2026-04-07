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

func setupDrawHandlerTestDB(t *testing.T) {
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

func TestDrawHandler_Draw_InsufficientPoints(t *testing.T) {
	setupDrawHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	prizeRepo := repository.NewPrizeRepository()
	drawService := service.NewDrawService(kidRepo, prizeRepo)
	handler := NewDrawHandler(drawService)

	r := gin.New()
	r.POST("/draw", handler.Draw)

	body, _ := json.Marshal(map[string]int{"kid_id": 1})
	req, _ := http.NewRequest("POST", "/draw", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDrawHandler_Draw_Success(t *testing.T) {
	setupDrawHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	prizeRepo := repository.NewPrizeRepository()
	drawService := service.NewDrawService(kidRepo, prizeRepo)
	handler := NewDrawHandler(drawService)

	kidRepo.UpdatePoints(1, 10)

	r := gin.New()
	r.POST("/draw", handler.Draw)

	body, _ := json.Marshal(map[string]int{"kid_id": 1})
	req, _ := http.NewRequest("POST", "/draw", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}
