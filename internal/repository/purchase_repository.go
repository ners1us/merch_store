package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
)

type PurchaseRepository interface {
	Create(purchase *model.Purchase) error
	GetUserPurchases(userID int) ([]model.InventoryItem, error)
}

type purchaseRepositoryImpl struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepositoryImpl{db: db}
}

func (pr *purchaseRepositoryImpl) Create(purchase *model.Purchase) error {
	return pr.db.Create(purchase).Error
}

func (pr *purchaseRepositoryImpl) GetUserPurchases(userID int) ([]model.InventoryItem, error) {
	var inventory []model.InventoryItem
	err := pr.db.Model(&model.Purchase{}).
		Select("merch_item as type, count(*) as quantity").
		Where("user_id = ?", userID).
		Group("merch_item").
		Scan(&inventory).Error
	return inventory, err
}
