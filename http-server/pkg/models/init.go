package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	BusOwner bool   `json:"busOwner"`
	BusID    *uint
	Bus      *Bus   `gorm:"foreignKey:BusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"  json:"bus,omitempty"`
}

type Bus struct {
	gorm.Model
	Name   string  `json:"name"`
	Owners []*User `json:"owners"`
}
