package service

import (
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService_GetUserInfo(t *testing.T) {
	// Arrange
	mockUserRepo := repository.NewMockUserRepository()
	mockPurchaseRepo := repository.NewMockPurchaseRepository()
	mockTransferRepo := repository.NewMockCoinTransferRepository()
	userService := NewUserService(mockUserRepo, mockPurchaseRepo, mockTransferRepo)

	user := &model.User{ID: 1, Coins: 1000}
	inventory := []model.InventoryItem{{Type: "socks", Quantity: 2}}
	received := []model.ReceivedCoinHistory{{FromUser: "alice", Amount: 100}}
	sent := []model.SentCoinHistory{{ToUser: "bob", Amount: 50}}

	mockUserRepo.On("FindByID", 1).Return(user, nil)
	mockPurchaseRepo.On("GetUserPurchases", 1).Return(inventory, nil)
	mockTransferRepo.On("GetReceivedTransfers", 1).Return(received, nil)
	mockTransferRepo.On("GetSentTransfers", 1).Return(sent, nil)

	// Act
	info, err := userService.GetUserInfo(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1000, info.Coins)
	assert.Equal(t, inventory, info.Inventory)
	assert.Equal(t, received, info.CoinHistory.Received)
	assert.Equal(t, sent, info.CoinHistory.Sent)

	// Arrange
	mockUserRepo.On("FindByID", 2).Return(&model.User{}, enum.ErrReceivingCoinsInfo)

	// Act
	info, err = userService.GetUserInfo(2)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Equal(t, enum.ErrReceivingCoinsInfo, err)
}
