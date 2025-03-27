package repository

import (
	"github.com/ners1us/merch_store/internal/model"
)

type PurchaseRepository interface {
	Create(purchase *model.Purchase) error
	GetUserPurchases(userID int) ([]model.InventoryItem, error)
}
