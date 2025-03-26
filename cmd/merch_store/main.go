package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/handler"
	"github.com/ners1us/merch_store/internal/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	dbConnection := os.Getenv("DB_URL")
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	var err error
	db, err := gorm.Open(postgres.Open(dbConnection), &gorm.Config{})
	if err != nil {
		return
	}

	handler.Init(db, jwtSecret)
	middleware.Init(jwtSecret)

	router := gin.Default()
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/auth", handler.HandleAuth)
		apiRoutes.GET("/info", middleware.AuthMiddleware(), handler.HandleInfo)
		apiRoutes.POST("/sendCoin", middleware.AuthMiddleware(), handler.HandleSendCoin)
		apiRoutes.GET("/buy/:item", middleware.AuthMiddleware(), handler.HandleBuy)
	}

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
