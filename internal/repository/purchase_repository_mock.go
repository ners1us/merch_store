package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockPurchaseRepository struct {
	mock.Mock
}

func NewMockPurchaseRepository() *MockPurchaseRepository {
	return &MockPurchaseRepository{}
}

func (mpr *MockPurchaseRepository) Create(purchase *model.Purchase) error {
	args := mpr.Called(purchase)
	return args.Error(0)
}

func (mpr *MockPurchaseRepository) GetUserPurchases(userID int) ([]model.InventoryItem, error) {
	args := mpr.Called(userID)
	return args.Get(0).([]model.InventoryItem), args.Error(1)
}
