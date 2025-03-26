package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/service"
	"net/http"
	"strconv"
)

type SendCoinHandler struct {
	transferService service.TransferService
}

func NewSendCoinHandler(transferService service.TransferService) *SendCoinHandler {
	return &SendCoinHandler{transferService: transferService}
}

func (sch *SendCoinHandler) HandleSendCoin(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrUserNotAuthorized.Error()})
		return
	}
	userID, _ := strconv.Atoi(userIDStr.(string))

	var req model.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrWrongReqFormat.Error()})
		return
	}

	err := sch.transferService.SendCoin(userID, req.ToUser, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": enum.SuccessfulTransfer.String()})
}
