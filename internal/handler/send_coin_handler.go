package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func HandleSendCoin(ctx *gin.Context) {
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrUserNotAuthorized.Error()})
		return
	}
	user := userInterface.(model.User)

	var request model.SendCoinRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrWrongReqFormat.Error()})
		return
	}
	if request.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrCoinsInappropriateAmount.Error()})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var sender model.User
		if err := tx.Where("id = ?", user.ID).
			First(&sender).Error; err != nil {
			return err
		}
		if sender.Coins < request.Amount {
			return enum.ErrInsufficientMoney
		}

		var receiver model.User
		if err := tx.Where("username = ?", request.ToUser).
			First(&receiver).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return enum.ErrReceiverNotFound
			}
			return err
		}

		if err := tx.Model(&sender).
			UpdateColumn("coins", gorm.Expr("coins - ?", request.Amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&receiver).
			UpdateColumn("coins", gorm.Expr("coins + ?", request.Amount)).Error; err != nil {
			return err
		}

		coinTransfer := model.CoinTransfer{
			FromUserID: sender.ID,
			ToUserID:   receiver.ID,
			Amount:     request.Amount,
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(&coinTransfer).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": enum.SuccessfulTransfer.String()})
}
