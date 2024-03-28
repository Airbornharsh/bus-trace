package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `gorm:"unique"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Phone    string `json:"phone" gorm:"unique"`
	BusOwner bool   `json:"busOwner"`
	BusID    *uint
	Bus      *Bus `gorm:"foreignKey:BusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"  json:"bus,omitempty"`
}

type Bus struct {
	gorm.Model
	ID   string  `gorm:"unique"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
