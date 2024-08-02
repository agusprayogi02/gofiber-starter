package entity

import "gorm.io/gorm"

type Post struct {
	ID     uint    `gorm:"primaryKey;autoIncrement"`
	Tweet  string  `gorm:"type:varchar(500)"`
	Photo  *string `gorm:"type:varchar(150)"`
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	gorm.Model
}
