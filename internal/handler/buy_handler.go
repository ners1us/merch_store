package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/service"
	"net/http"
	"strconv"
)

type BuyHandler struct {
	merchService service.MerchService
}

func NewBuyHandler(merchService service.MerchService) *BuyHandler {
	return &BuyHandler{merchService: merchService}
}

func (bh *BuyHandler) HandleBuy(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrUserNotAuthorized.Error()})
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrNotProvidedItem.Error()})
		return
	}

	err = bh.merchService.BuyMerch(userID, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": enum.SuccessfulPurchase.String()})
}
