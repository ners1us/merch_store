package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
)

type coinTransferRepository struct {
	db *gorm.DB
}

func NewCoinTransferRepository(db *gorm.DB) CoinTransferRepository {
	return &coinTransferRepository{db: db}
}

func (ctr *coinTransferRepository) Create(transfer *model.CoinTransfer) error {
	return ctr.db.Create(transfer).Error
}

func (ctr *coinTransferRepository) GetReceivedTransfers(userID int) ([]model.ReceivedCoinHistory, error) {
	var received []model.ReceivedCoinHistory
	err := ctr.db.Table("coin_transfers").
		Select("users.username as from_user, coin_transfers.amount").
		Joins("join users on coin_transfers.from_user_id = users.id").
		Where("coin_transfers.to_user_id = ?", userID).
		Scan(&received).Error
	return received, err
}

func (ctr *coinTransferRepository) GetSentTransfers(userID int) ([]model.SentCoinHistory, error) {
	var sent []model.SentCoinHistory
	err := ctr.db.Table("coin_transfers").
		Select("users.username as to_user, coin_transfers.amount").
		Joins("join users on coin_transfers.to_user_id = users.id").
		Where("coin_transfers.from_user_id = ?", userID).
		Scan(&sent).Error
	return sent, err
}
