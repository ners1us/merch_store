package repository

import (
	"github.com/ners1us/merch_store/internal/model"
	"gorm.io/gorm"
)

type coinTransferRepositoryImpl struct {
	db *gorm.DB
}

func NewCoinTransferRepository(db *gorm.DB) CoinTransferRepository {
	return &coinTransferRepositoryImpl{db: db}
}

func (ctr *coinTransferRepositoryImpl) Create(transfer *model.CoinTransfer) error {
	return ctr.db.Create(transfer).Error
}

func (ctr *coinTransferRepositoryImpl) GetReceivedTransfers(userID int) ([]model.ReceivedCoinHistory, error) {
	var received []model.ReceivedCoinHistory
	err := ctr.db.Table("coin_transfers").
		Select("users.username as from_user, coin_transfers.amount").
		Joins("join users on coin_transfers.from_user_id = users.id").
		Where("coin_transfers.to_user_id = ?", userID).
		Scan(&received).Error
	return received, err
}

func (ctr *coinTransferRepositoryImpl) GetSentTransfers(userID int) ([]model.SentCoinHistory, error) {
	var sent []model.SentCoinHistory
	err := ctr.db.Table("coin_transfers").
		Select("users.username as to_user, coin_transfers.amount").
		Joins("join users on coin_transfers.to_user_id = users.id").
		Where("coin_transfers.from_user_id = ?", userID).
		Scan(&sent).Error
	return sent, err
}
