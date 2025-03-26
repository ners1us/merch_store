package service

import (
	"errors"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
	"gorm.io/gorm"
	"time"
)

type transferService struct {
	userRepo     repository.UserRepository
	transferRepo repository.CoinTransferRepository
}

func NewTransferService(userRepo repository.UserRepository, transferRepo repository.CoinTransferRepository) TransferService {
	return &transferService{userRepo: userRepo, transferRepo: transferRepo}
}

func (ts *transferService) SendCoin(fromUserID int, toUsername string, amount int) error {
	if amount <= 0 {
		return enum.ErrCoinsInappropriateAmount
	}

	return ts.userRepo.RunTransaction(func(tx *gorm.DB) error {
		sender, err := ts.userRepo.FindByID(fromUserID)
		if err != nil {
			return err
		}

		if sender.Coins < amount {
			return enum.ErrInsufficientMoney
		}

		receiver, err := ts.userRepo.FindByUsername(toUsername)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return enum.ErrReceiverNotFound
			}
			return err
		}

		if sender.Username == receiver.Username {
			return enum.ErrEqualReceivers
		}

		sender.Coins -= amount
		receiver.Coins += amount

		if err := ts.userRepo.Update(sender); err != nil {
			return err
		}
		if err := ts.userRepo.Update(receiver); err != nil {
			return err
		}

		transfer := &model.CoinTransfer{
			FromUserID: sender.ID,
			ToUserID:   receiver.ID,
			Amount:     amount,
			CreatedAt:  time.Now(),
		}
		if err := ts.transferRepo.Create(transfer); err != nil {
			return err
		}

		return nil
	})
}
