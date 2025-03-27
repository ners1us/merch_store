package service

type TransferService interface {
	SendCoin(fromUserID int, toUsername string, amount int) error
}
