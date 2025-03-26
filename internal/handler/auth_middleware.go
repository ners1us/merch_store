package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrNoAuthToken.Error()})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrWrongTokenFormat.Error()})
			c.Abort()
			return
		}
		tokenStr := parts[1]

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": enum.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", strconv.Itoa(claims.UserID))
		c.Set("username", claims.Username)
		c.Next()
	}
}
