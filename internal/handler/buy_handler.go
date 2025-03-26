package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enums"
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func HandleBuy(ctx *gin.Context) {
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": enums.ErrUserNotAuthorized.Error()})
		return
	}
	user := userInterface.(model.User)

	item := ctx.Param("item")
	if item == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": enums.ErrNotProvidedItem.Error()})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var merch model.Merch
		if err := tx.Where("name = ?", item).First(&merch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return enums.ErrItemNotFound
			}
			return err
		}

		var currentUser model.User
		if err := tx.Where("id = ?", user.ID).
			First(&currentUser).Error; err != nil {
			return err
		}
		if currentUser.Coins < merch.Price {
			return enums.ErrBuyWithInsufficientMoney
		}

		if err := tx.Model(&currentUser).
			UpdateColumn("coins", gorm.Expr("coins - ?", merch.Price)).Error; err != nil {
			return err
		}

		purchase := model.Purchase{
			UserID:    currentUser.ID,
			MerchItem: item,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": enums.SuccessfulPurchase.String()})
}
