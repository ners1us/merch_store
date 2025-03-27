package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (mur *MockUserRepository) Create(user *model.User) error {
	args := mur.Called(user)
	return args.Error(0)
}

func (mur *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	args := mur.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (mur *MockUserRepository) FindByID(id int) (*model.User, error) {
	args := mur.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (mur *MockUserRepository) Update(user *model.User) error {
	args := mur.Called(user)
	return args.Error(0)
}

func (mur *MockUserRepository) RunTransaction(fn func(tx *gorm.DB) error) error {
	args := mur.Called(fn)
	return args.Error(0)
}
