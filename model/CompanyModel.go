package model

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Id      uint   `json:"id" gorm:"primaryKey;autoIncrement"` // Id is the primary key
	UserID  uint   `json:"user_id" gorm:"not null"`            // UserID is the foreign key
	Name    string `json:"name" gorm:"size:255;unique;not null"`
	Address string `json:"address" gorm:"size:255;not null"`
}
