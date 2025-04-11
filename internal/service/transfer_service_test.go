package service

import (
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestTransferService_SendCoin(t *testing.T) {
	// Arrange
	mockUserRepo := repository.NewMockUserRepository()
	mockTransferRepo := repository.NewMockCoinTransferRepository()
	transferService := NewTransferService(mockUserRepo, mockTransferRepo)

	sender := &model.User{ID: 1, Username: "alice", Coins: 1000}
	receiver := &model.User{ID: 2, Username: "bob", Coins: 500}

	mockUserRepo.On("FindByID", 1).Return(sender, nil).Once()
	mockUserRepo.On("FindByUsername", "bob").Return(receiver, nil).Once()
	mockUserRepo.On("Update", mock.Anything).Return(nil).Times(2)
	mockTransferRepo.On("Create", mock.Anything).Return(nil).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.NoError(t, err)
	}).Return(nil).Once()

	// Act
	err := transferService.SendCoin(1, "bob", 200)

	// Assert
	assert.NoError(t, err)

	// Arrange
	sender.Coins = 100
	mockUserRepo.On("FindByID", 1).Return(sender, nil).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.Equal(t, enum.ErrInsufficientMoney, err)
	}).Return(enum.ErrInsufficientMoney).Once()

	// Act
	err = transferService.SendCoin(1, "bob", 200)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrInsufficientMoney, err)

	// Arrange
	mockUserRepo.On("FindByID", 1).Return(sender, nil).Once()
	mockUserRepo.On("FindByUsername", "alice").Return(sender, nil).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.Equal(t, enum.ErrEqualReceivers, err)
	}).Return(enum.ErrEqualReceivers).Once()

	// Act
	err = transferService.SendCoin(1, "alice", 100)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrEqualReceivers, err)

	// Act
	err = transferService.SendCoin(1, "bob", -10)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrCoinsInappropriateAmount, err)
}
