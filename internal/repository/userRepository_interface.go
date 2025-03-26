package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id int) (*model.User, error)
	Update(user *model.User) error
	RunTransaction(fn func(tx *gorm.DB) error) error
}
