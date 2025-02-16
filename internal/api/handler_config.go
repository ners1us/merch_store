package api

import (
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	jwtSecret []byte
)

func Init(database *gorm.DB, secret []byte) {
	db = database
	jwtSecret = secret
}
