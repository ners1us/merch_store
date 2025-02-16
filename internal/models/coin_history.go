package models

type CoinHistory struct {
	Received []ReceivedCoinHistory `json:"received"`
	Sent     []SentCoinHistory     `json:"sent"`
}
