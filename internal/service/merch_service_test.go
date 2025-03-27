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

func TestMerchService_BuyMerch(t *testing.T) {
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

	err := merchService.BuyMerch(1, "pink-hoody")
	assert.NoError(t, err)

	user.Coins = 400
	mockMerchRepo.On("FindByName", "pink-hoody").Return(merch, nil).Once()
	mockUserRepo.On("FindByID", 1).Return(user, nil).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.Equal(t, enum.ErrBuyWithInsufficientMoney, err)
	}).Return(enum.ErrBuyWithInsufficientMoney).Once()

	err = merchService.BuyMerch(1, "pink-hoody")
	assert.Error(t, err)
	assert.Equal(t, enum.ErrBuyWithInsufficientMoney, err)

	mockMerchRepo.On("FindByName", "candy").Return(&model.Merch{}, gorm.ErrRecordNotFound).Once()
	mockUserRepo.On("RunTransaction", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(tx *gorm.DB) error)
		err := fn(nil)
		assert.Equal(t, enum.ErrItemNotFound, err)
	}).Return(enum.ErrItemNotFound).Once()

	err = merchService.BuyMerch(1, "candy")
	assert.Error(t, err)
	assert.Equal(t, enum.ErrItemNotFound, err)
}
