package main

import (
	"github.com/ners1us/merch_store/internal/model"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/config"
	"github.com/ners1us/merch_store/internal/handler"
	"github.com/ners1us/merch_store/internal/repository"
	"github.com/ners1us/merch_store/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.InitConfig()

	db, err := gorm.Open(postgres.Open(cfg.DbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Merch{}, &model.Purchase{}, &model.CoinTransfer{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	userRepo := repository.NewUserRepository(db)
	merchRepo := repository.NewMerchRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	transferRepo := repository.NewCoinTransferRepository(db)

	err = merchRepo.InitializeMerch()
	if err != nil {
		log.Fatal("failed to initialize data in the merch table: ", err)
	}

	authService := service.NewAuthService(userRepo, []byte(cfg.JWTSecret))
	userService := service.NewUserService(userRepo, purchaseRepo, transferRepo)
	merchService := service.NewMerchService(userRepo, merchRepo, purchaseRepo)
	transferService := service.NewTransferService(userRepo, transferRepo)

	authHandler := handler.NewAuthHandler(authService)
	infoHandler := handler.NewInfoHandler(userService)
	buyHandler := handler.NewBuyHandler(merchService)
	sendCoinHandler := handler.NewSendCoinHandler(transferService)

	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/auth", authHandler.HandleAuth)
		api.GET("/info", handler.AuthMiddleware([]byte(cfg.JWTSecret)), infoHandler.HandleInfo)
		api.POST("/sendCoin", handler.AuthMiddleware([]byte(cfg.JWTSecret)), sendCoinHandler.HandleSendCoin)
		api.GET("/buy/:item", handler.AuthMiddleware([]byte(cfg.JWTSecret)), buyHandler.HandleBuy)
	}

	err = r.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("failed running merch store service: ", err)
	}
}
