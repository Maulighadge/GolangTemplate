package models

import (
	"github.com/google/uuid"
)

type Product struct {
	BaseModel
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	SKU         string    `gorm:"size:100;uniqueIndex;not null" json:"sku"`
	Price       float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"price"`
	Stock       int       `gorm:"default:0;not null" json:"stock"`
	Category    string    `gorm:"size:100;index" json:"category"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
