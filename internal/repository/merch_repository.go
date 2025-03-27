package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
)

type merchRepository struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) MerchRepository {
	return &merchRepository{db: db}
}

func (mr *merchRepository) FindByName(name string) (*model.Merch, error) {
	var merch model.Merch
	err := mr.db.Where("name = ?", name).First(&merch).Error
	return &merch, err
}

func (mr *merchRepository) InitializeMerch() error {
	merch := []model.Merch{
		{"t-shirt", 20},
		{"cup", 20},
		{"book", 50},
		{"pen", 10},
		{"powerbank", 200},
		{"hoody", 300},
		{"umbrella", 200},
		{"socks", 10},
		{"wallet", 50},
		{"pink-hoody", 500}}
	return mr.db.Create(&merch).Error
}
