package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"net/http"
)

func HandleInfo(ctx *gin.Context) {
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrUserNotAuthorized.Error()})
		return
	}
	user := userInterface.(model.User)

	var dbUser model.User
	if err := db.First(&dbUser, user.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrReceivingCoinsInfo.Error()})
		return
	}
	coins := dbUser.Coins

	var inventory []model.InventoryItem
	if err := db.Model(&model.Purchase{}).
		Select("merch_item as type, count(*) as quantity").
		Where("user_id = ?", user.ID).
		Group("merch_item").
		Scan(&inventory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrReceivingPurchaseHistory.Error()})
		return
	}

	var received []model.ReceivedCoinHistory
	if err := db.Table("coin_transfers").
		Select("users.username as from_user, coin_transfers.amount").
		Joins("join users on coin_transfers.from_user_id = users.id").
		Where("coin_transfers.to_user_id = ?", user.ID).
		Scan(&received).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrReceivingTransferHistory.Error()})
		return
	}

	var sent []model.SentCoinHistory
	if err := db.Table("coin_transfers").
		Select("users.username as to_user, coin_transfers.amount").
		Joins("join users on coin_transfers.to_user_id = users.id").
		Where("coin_transfers.from_user_id = ?", user.ID).
		Scan(&sent).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrReceivingTransferHistory.Error()})
		return
	}

	response := model.InfoResponse{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: model.CoinHistory{Received: received, Sent: sent},
	}
	ctx.JSON(http.StatusOK, response)
}
