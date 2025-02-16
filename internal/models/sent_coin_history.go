package models

type SentCoinHistory struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
