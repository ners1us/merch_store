package service

import (
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

func TestAuthService_Authenticate(t *testing.T) {
	// Arrange
	mockUserRepo := repository.NewMockUserRepository()
	jwtSecret := []byte("secret")
	authService := NewAuthService(mockUserRepo, jwtSecret)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("cool_password"), bcrypt.DefaultCost)
	existingUser := &model.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
	}
	mockUserRepo.On("FindByUsername", "testuser").Return(existingUser, nil)

	// Act
	token, err := authService.Authenticate("testuser", "cool_password")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Arrange
	mockUserRepo.On("FindByUsername", "newuser").Return(&model.User{}, gorm.ErrRecordNotFound)
	mockUserRepo.On("Create", mock.Anything).Return(nil)

	// Act
	token, err = authService.Authenticate("newuser", "new_password")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Arrange
	mockUserRepo.On("FindByUsername", "testuser").Return(existingUser, nil)

	// Act
	token, err = authService.Authenticate("testuser", "wrong_password")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrWrongCredentials, err)
	assert.Empty(t, token)
}
