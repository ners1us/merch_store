package model

type ReceivedCoinHistory struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}
