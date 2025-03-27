package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockMerchRepository struct {
	mock.Mock
}

func NewMockMerchRepository() *MockMerchRepository {
	return &MockMerchRepository{}
}

func (mmr *MockMerchRepository) FindByName(name string) (*model.Merch, error) {
	args := mmr.Called(name)
	return args.Get(0).(*model.Merch), args.Error(1)
}

// InitializeMerch Not implemented
func (mmr *MockMerchRepository) InitializeMerch() error {
	return nil
}
