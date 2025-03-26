package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func HandleAuth(ctx *gin.Context) {
	var request model.AuthRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrWrongReqFormat.Error()})
		return
	}
	if request.Username == "" || request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrNoUsernameAndPassword.Error()})
		return
	}

	var user model.User
	result := db.Where("username = ?", request.Username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrCreatingUser.Error()})
			return
		}
		user = model.User{
			Username: request.Username,
			Password: string(hash),
			Coins:    1000,
		}
		if err := db.Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrCreatingUser.Error()})
			return
		}
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrInternalServer.Error()})
		return
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrWrongCredentials.Error()})
			return
		}
	}

	expTime := time.Now().Add(time.Hour)
	claims := &model.Claims{
		Username: user.Username,
		UserID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := tokenObj.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": enum.ErrGeneratingToken.Error()})
		return
	}

	ctx.JSON(http.StatusOK, model.AuthResponse{Token: tokenStr})
}
