package service

import (
	"github.com/ners1us/merch_store/internal/model"
)

type UserService interface {
	GetUserInfo(userID int) (*model.InfoResponse, error)
}
