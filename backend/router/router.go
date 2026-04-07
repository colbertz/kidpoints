package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"points-app/database"
	"points-app/handlers"
	"points-app/repository"
	"points-app/service"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	kidRepo := repository.NewKidRepository()
	behaviorRepo := repository.NewBehaviorRepository()
	prizeRepo := repository.NewPrizeRepository()
	recordRepo := repository.NewRecordRepository()

	kidService := service.NewKidService(kidRepo, behaviorRepo, recordRepo)
	drawService := service.NewDrawService(kidRepo, prizeRepo)
	recordService := service.NewRecordService(recordRepo, kidRepo)

	kidHandler := handlers.NewKidHandler(kidService)
	behaviorHandler := handlers.NewBehaviorHandler(behaviorRepo)
	prizeHandler := handlers.NewPrizeHandler(prizeRepo)
	drawHandler := handlers.NewDrawHandler(drawService)
	recordHandler := handlers.NewRecordHandler(recordService)

	r.GET("/api/kids", kidHandler.GetKids)
	r.POST("/api/kids", kidHandler.CreateKid)
	r.GET("/api/kids/:id/points", kidHandler.GetPoints)
	r.PUT("/api/kids/:id", kidHandler.UpdateKid)
	r.DELETE("/api/kids/:id", kidHandler.DeleteKid)
	r.POST("/api/kids/points/add", kidHandler.AddPoints)
	r.POST("/api/kids/points/sub", kidHandler.SubPoints)

	r.GET("/api/behaviors", behaviorHandler.GetBehaviors)
	r.POST("/api/behaviors", behaviorHandler.CreateBehavior)
	r.PUT("/api/behaviors/:id", behaviorHandler.UpdateBehavior)
	r.DELETE("/api/behaviors/:id", behaviorHandler.DeleteBehavior)

	r.GET("/api/prizes", prizeHandler.GetPrizes)
	r.POST("/api/prizes", prizeHandler.CreatePrize)
	r.PUT("/api/prizes/:id", prizeHandler.UpdatePrize)
	r.DELETE("/api/prizes/:id", prizeHandler.DeletePrize)

	r.POST("/api/draw", drawHandler.Draw)

	r.GET("/api/records", recordHandler.GetRecords)
	r.POST("/api/records/:id/reverse", recordHandler.ReverseRecord)

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "db": database.DB != nil})
	})

	return r
}
