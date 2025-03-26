package service

type MerchService interface {
	BuyMerch(userID int, item string) error
}
