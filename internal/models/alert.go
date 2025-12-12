package models

type Alert struct {
	ID        uint   `gorm:"primaryKey"`
	Symbol    string `gorm:"size:255;not null"`
	Condition string `gorm:"size:1;not null"`
	Price     float64    `gorm:"not null"`
	Triggered bool   `gorm:"default:false"`
}
