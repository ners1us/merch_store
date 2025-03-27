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
