package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/service"
	"net/http"
	"strconv"
)

type InfoHandler struct {
	userService service.UserService
}

func NewInfoHandler(userService service.UserService) *InfoHandler {
	return &InfoHandler{userService: userService}
}

func (ih *InfoHandler) HandleInfo(c *gin.Context) {
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

	info, err := ih.userService.GetUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
