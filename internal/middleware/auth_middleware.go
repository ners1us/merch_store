package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ners1us/merch_store/internal/enums"
	"github.com/ners1us/merch_store/internal/models"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": enums.ErrNoAuthToken.Error()})
			ctx.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if parts[0] != "Bearer" || len(parts) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": enums.ErrWrongTokenFormat.Error()})
			ctx.Abort()
			return
		}
		tokenStr := parts[1]

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": enums.ErrInvalidToken.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("user", models.User{
			ID:       claims.UserID,
			Username: claims.Username,
		})
		ctx.Next()
	}
}
