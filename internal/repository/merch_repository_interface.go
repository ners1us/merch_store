package repository

import "github.com/ners1us/merch_store/internal/model"

type MerchRepository interface {
	FindByName(name string) (*model.Merch, error)
	InitializeMerch() error
}
