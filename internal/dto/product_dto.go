package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/internal/models"
)

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Description string  `json:"description" binding:"omitempty"`
	SKU         string  `json:"sku" binding:"required,min=2,max=100"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Category    string  `json:"category" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string  `json:"description" binding:"omitempty"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	Category    string   `json:"category" binding:"omitempty"`
}

type ProductResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	SKU         string        `json:"sku"`
	Price       float64       `json:"price"`
	Stock       int           `json:"stock"`
	Category    string        `json:"category"`
	UserID      uuid.UUID     `json:"user_id"`
	User        *UserResponse `json:"user,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func ToProductResponse(product *models.Product) ProductResponse {
	resp := ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		UserID:      product.UserID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	if product.User != nil {
		u := ToUserResponse(product.User)
		resp.User = &u
	}

	return resp
}

func ToProductResponseList(products []models.Product) []ProductResponse {
	list := make([]ProductResponse, len(products))
	for i, p := range products {
		list[i] = ToProductResponse(&p)
	}
	return list
}
