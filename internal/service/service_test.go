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

func TestMerchService_BuyMerch(t *testing.T) {
	// Arrange
	mockUserRepo := repository.NewMockUserRepository()
	mockMerchRepo := repository.NewMockMerchRepository()
	mockPurchaseRepo := repository.NewMockPurchaseRepository()
	merchService := NewMerchService(mockUserRepo, mockMerchRepo, mockPurchaseRepo)

	user := &model.User{ID: 1, Coins: 1000}
	merch := &model.Merch{Name: "pink-hoody", Price: 500}

	mockMerchRepo.On("FindByName", "pink-hoody").Return(merch, nil).Once()
	mockUserRepo.On("FindByID", 1).Return(user, nil).Once()
	mockUserRepo.On("Update", mock.Anything).Return(nil).Once()
	mockPurchaseRepo.On("Create", mock.Anything).Return(nil).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.NoError(t, err)
	}).Return(nil).Once()

	// Act
	err := merchService.BuyMerch(1, "pink-hoody")

	// Assert
	assert.NoError(t, err)

	// Arrange
	user.Coins = 400
	mockMerchRepo.On("FindByName", "pink-hoody").Return(merch, nil).Once()
	mockUserRepo.On("FindByID", 1).Return(user, nil).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.Equal(t, enum.ErrBuyWithInsufficientMoney, err)
	}).Return(enum.ErrBuyWithInsufficientMoney).Once()

	// Act
	err = merchService.BuyMerch(1, "pink-hoody")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrBuyWithInsufficientMoney, err)

	// Arrange
	mockMerchRepo.On("FindByName", "candy").Return(&model.Merch{}, gorm.ErrRecordNotFound).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.Equal(t, enum.ErrItemNotFound, err)
	}).Return(enum.ErrItemNotFound).Once()

	// Act
	err = merchService.BuyMerch(1, "candy")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrItemNotFound, err)
}

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
