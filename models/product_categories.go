package models

type ProductCategory struct {
	ID   uint   `gorm:"primaryKey" json:"-"`
	Code string `gorm:"uniqueIndex;not null" json:"code"`
	Name string `json:"name"`
}
