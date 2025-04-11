package service

import (
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
)

type UserService interface {
	GetUserInfo(userID int) (*model.InfoResponse, error)
}

type userServiceImpl struct {
	userRepo     repository.UserRepository
	purchaseRepo repository.PurchaseRepository
	transferRepo repository.CoinTransferRepository
}

func NewUserService(userRepo repository.UserRepository, purchaseRepo repository.PurchaseRepository, transferRepo repository.CoinTransferRepository) UserService {
	return &userServiceImpl{userRepo: userRepo, purchaseRepo: purchaseRepo, transferRepo: transferRepo}
}

func (us *userServiceImpl) GetUserInfo(userID int) (*model.InfoResponse, error) {
	user, err := us.userRepo.FindByID(userID)
	if err != nil {
		return nil, enum.ErrReceivingCoinsInfo
	}

	inventory, err := us.purchaseRepo.GetUserPurchases(userID)
	if err != nil {
		return nil, enum.ErrReceivingPurchaseHistory
	}

	received, err := us.transferRepo.GetReceivedTransfers(userID)
	if err != nil {
		return nil, enum.ErrReceivingTransferHistory
	}
	sent, err := us.transferRepo.GetSentTransfers(userID)
	if err != nil {
		return nil, enum.ErrReceivingTransferHistory
	}

	return &model.InfoResponse{
		Coins:     user.Coins,
		Inventory: inventory,
		CoinHistory: model.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
