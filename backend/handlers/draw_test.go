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
	// 使用 init.go 中的奖项数据
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

func TestDrawHandler_DrawProbabilityDistribution(t *testing.T) {
	setupDrawHandlerTestDB(t)
	kidRepo := repository.NewKidRepository()
	prizeRepo := repository.NewPrizeRepository()
	drawService := service.NewDrawService(kidRepo, prizeRepo)
	handler := NewDrawHandler(drawService)

	// 获取所有奖项
	prizes, _ := prizeRepo.GetAll()
	prizeCount := make(map[int]int)
	totalDraws := 1000

	// 给 kid 足够积分
	kidRepo.UpdatePoints(1, totalDraws*2+100)

	r := gin.New()
	r.POST("/draw", handler.Draw)

	for i := 0; i < totalDraws; i++ {
		body, _ := json.Marshal(map[string]int{"kid_id": 1})
		req, _ := http.NewRequest("POST", "/draw", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			var result models.Prize
			json.Unmarshal(w.Body.Bytes(), &result)
			prizeCount[result.ID]++
		}
	}

	// 打印统计结果
	t.Log("========== 抽奖概率统计 ==========")
	t.Logf("总抽奖次数: %d", totalDraws)
	t.Log("")
	for _, p := range prizes {
		count := prizeCount[p.ID]
		actualRate := float64(count) / float64(totalDraws) * 100
		t.Logf("%s: 期望 %.1f%%, 实际 %.1f%% (抽中 %d次)", p.Name, p.Probability, actualRate, count)
	}
	t.Log("==================================")
}
