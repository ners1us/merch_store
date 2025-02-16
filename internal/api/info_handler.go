package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enums"
	"github.com/ners1us/merch_store/internal/models"
	"net/http"
)

func HandleInfo(ctx *gin.Context) {
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": enums.ErrUserNotAuthorized.Error()})
		return
	}
	user := userInterface.(models.User)

	var dbUser models.User
	if err := db.First(&dbUser, user.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enums.ErrReceivingCoinsInfo.Error()})
		return
	}
	coins := dbUser.Coins

	var inventory []models.InventoryItem
	if err := db.Model(&models.Purchase{}).
		Select("merch_item as type, count(*) as quantity").
		Where("user_id = ?", user.ID).
		Group("merch_item").
		Scan(&inventory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enums.ErrReceivingPurchaseHistory.Error()})
		return
	}

	var received []models.ReceivedCoinHistory
	if err := db.Table("coin_transfers").
		Select("users.username as from_user, coin_transfers.amount").
		Joins("join users on coin_transfers.from_user_id = users.id").
		Where("coin_transfers.to_user_id = ?", user.ID).
		Scan(&received).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enums.ErrReceivingTransferHistory.Error()})
		return
	}

	var sent []models.SentCoinHistory
	if err := db.Table("coin_transfers").
		Select("users.username as to_user, coin_transfers.amount").
		Joins("join users on coin_transfers.to_user_id = users.id").
		Where("coin_transfers.from_user_id = ?", user.ID).
		Scan(&sent).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enums.ErrReceivingTransferHistory.Error()})
		return
	}

	response := models.InfoResponse{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: models.CoinHistory{Received: received, Sent: sent},
	}
	ctx.JSON(http.StatusOK, response)
}
