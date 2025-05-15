package entity

import "gorm.io/gorm"

type ProcessedApartment struct {
	gorm.Model
}

func (*ProcessedApartment) TableName() string {
	return "processed_apartments"
}
