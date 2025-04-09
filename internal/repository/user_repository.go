package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (ur *userRepositoryImpl) Create(user *model.User) error {
	return ur.db.Create(user).Error
}

func (ur *userRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (ur *userRepositoryImpl) FindByID(id int) (*model.User, error) {
	var user model.User
	err := ur.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (ur *userRepositoryImpl) Update(user *model.User) error {
	return ur.db.Save(user).Error
}

func (ur *userRepositoryImpl) RunTransaction(fn func(tx *gorm.DB) error) error {
	return ur.db.Transaction(fn)
}
