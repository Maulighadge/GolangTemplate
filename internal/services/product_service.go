package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/internal/dto"
	"github.com/Maulighadge/GolangTemplate/internal/models"
	"github.com/Maulighadge/GolangTemplate/internal/repositories"
	"github.com/Maulighadge/GolangTemplate/pkg/utils"
)

type ProductService interface {
	Create(userID uuid.UUID, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetByID(id uuid.UUID) (*dto.ProductResponse, error)
	List(search, category string, page, limit int) ([]dto.ProductResponse, utils.Pagination, error)
	Update(id uuid.UUID, userID uuid.UUID, userRole models.Role, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Delete(id uuid.UUID, userID uuid.UUID, userRole models.Role) error
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) Create(userID uuid.UUID, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Verify unique SKU
	existing, err := s.productRepo.FindBySKU(req.SKU)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		UserID:      userID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	res := dto.ToProductResponse(product)
	return &res, nil
}

func (s *productService) GetByID(id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	res := dto.ToProductResponse(product)
	return &res, nil
}

func (s *productService) List(search, category string, page, limit int) ([]dto.ProductResponse, utils.Pagination, error) {
	pagination := utils.Pagination{
		Page:  page,
		Limit: limit,
	}

	products, totalRows, err := s.productRepo.FindAll(search, category, pagination)
	if err != nil {
		return nil, pagination, err
	}

	meta := utils.CalculatePagination(page, limit, totalRows)
	respList := dto.ToProductResponseList(products)

	return respList, meta, nil
}

func (s *productService) Update(id uuid.UUID, userID uuid.UUID, userRole models.Role, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Authorization check: Only owner or admin can update
	if userRole != models.RoleAdmin && product.UserID != userID {
		return nil, errors.New("you do not have permission to update this product")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	res := dto.ToProductResponse(product)
	return &res, nil
}

func (s *productService) Delete(id uuid.UUID, userID uuid.UUID, userRole models.Role) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	// Authorization check: Only owner or admin can delete
	if userRole != models.RoleAdmin && product.UserID != userID {
		return errors.New("you do not have permission to delete this product")
	}

	return s.productRepo.Delete(id)
}
