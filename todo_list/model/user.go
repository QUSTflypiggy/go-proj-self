package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	PasswordDigest string //存储的是密文，加密后的密码
}
