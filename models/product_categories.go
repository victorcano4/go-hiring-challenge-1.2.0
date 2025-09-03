package models

type ProductCategory struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string
}
