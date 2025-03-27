package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockCoinTransferRepository struct {
	mock.Mock
}

func NewMockCoinTransferRepository() *MockCoinTransferRepository {
	return &MockCoinTransferRepository{}
}

func (mctr *MockCoinTransferRepository) Create(transfer *model.CoinTransfer) error {
	args := mctr.Called(transfer)
	return args.Error(0)
}

func (mctr *MockCoinTransferRepository) GetReceivedTransfers(userID int) ([]model.ReceivedCoinHistory, error) {
	args := mctr.Called(userID)
	return args.Get(0).([]model.ReceivedCoinHistory), args.Error(1)
}

func (mctr *MockCoinTransferRepository) GetSentTransfers(userID int) ([]model.SentCoinHistory, error) {
	args := mctr.Called(userID)
	return args.Get(0).([]model.SentCoinHistory), args.Error(1)
}
