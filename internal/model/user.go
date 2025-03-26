package model

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Coins    int    `gorm:"not null;default:1000" json:"coins"`
}
