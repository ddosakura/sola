package main

import (
	"github.com/jinzhu/gorm"
)

// User Model
type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Secret   string `gorm:"not null"`
}

func init() {
	db.AutoMigrate(&User{})
}
