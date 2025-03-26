package model

type Merch struct {
	Name  string `gorm:"primaryKey;not null" json:"name"`
	Price int    `gorm:"not null" json:"price"`
}
