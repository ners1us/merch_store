package repository

import "github.com/ners1us/merch_store/internal/model"

type CoinTransferRepository interface {
	Create(transfer *model.CoinTransfer) error
	GetReceivedTransfers(userID int) ([]model.ReceivedCoinHistory, error)
	GetSentTransfers(userID int) ([]model.SentCoinHistory, error)
}
