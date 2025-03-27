package service

import (
	"errors"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"

	"github.com/golang-jwt/jwt"
)

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret []byte) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (as *authService) Authenticate(username, password string) (string, error) {
	user, err := as.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return "", enum.ErrCreatingUser
			}
			user = &model.User{
				Username: username,
				Password: string(hash),
				Coins:    1000,
			}
			if err := as.userRepo.Create(user); err != nil {
				return "", enum.ErrCreatingUser
			}
		} else {
			return "", enum.ErrInternalServer
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return "", enum.ErrWrongCredentials
		}
	}

	claims := &model.Claims{
		Username: user.Username,
		UserID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(as.jwtSecret)
	if err != nil {
		return "", enum.ErrGeneratingToken
	}
	return tokenStr, nil
}
