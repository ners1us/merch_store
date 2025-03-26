package model

import "time"

type CoinTransfer struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	FromUserID int       `gorm:"not null" json:"from_user_id"`
	ToUserID   int       `gorm:"not null" json:"to_user_id"`
	Amount     int       `gorm:"not null" json:"amount"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}
