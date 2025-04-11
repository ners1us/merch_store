package service

import (
	"errors"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
	"gorm.io/gorm"
	"time"
)

type MerchService interface {
	BuyMerch(userID int, item string) error
}

type merchServiceImpl struct {
	userRepo     repository.UserRepository
	merchRepo    repository.MerchRepository
	purchaseRepo repository.PurchaseRepository
}

func NewMerchService(userRepo repository.UserRepository, merchRepo repository.MerchRepository, purchaseRepo repository.PurchaseRepository) MerchService {
	return &merchServiceImpl{userRepo: userRepo, merchRepo: merchRepo, purchaseRepo: purchaseRepo}
}

func (ms *merchServiceImpl) BuyMerch(userID int, item string) error {
	return ms.userRepo.RunTransaction(func(tx *gorm.DB) error {
		merch, err := ms.merchRepo.FindByName(item)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return enum.ErrItemNotFound
			}
			return err
		}

		user, err := ms.userRepo.FindByID(userID)
		if err != nil {
			return err
		}

		if user.Coins < merch.Price {
			return enum.ErrBuyWithInsufficientMoney
		}

		user.Coins -= merch.Price
		if err := ms.userRepo.Update(user); err != nil {
			return err
		}

		purchase := &model.Purchase{
			UserID:    userID,
			MerchItem: item,
			CreatedAt: time.Now(),
		}
		if err := ms.purchaseRepo.Create(purchase); err != nil {
			return err
		}

		return nil
	})
}
