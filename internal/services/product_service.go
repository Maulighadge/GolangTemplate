package services

import (
	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/internal/dto"
	"github.com/Maulighadge/GolangTemplate/internal/models"
	"github.com/Maulighadge/GolangTemplate/internal/repositories"
	"github.com/Maulighadge/GolangTemplate/pkg/apperror"
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
		return nil, apperror.NewInternal(err, "failed to query product by SKU")
	}
	if existing != nil {
		return nil, apperror.ErrSKUAlreadyExists
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
		return nil, apperror.NewInternal(err, "failed to create product")
	}

	res := dto.ToProductResponse(product)
	return &res, nil
}

func (s *productService) GetByID(id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to retrieve product")
	}
	if product == nil {
		return nil, apperror.ErrProductNotFound
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
		return nil, pagination, apperror.NewInternal(err, "failed to list products")
	}

	meta := utils.CalculatePagination(page, limit, totalRows)
	respList := dto.ToProductResponseList(products)

	return respList, meta, nil
}

func (s *productService) Update(id uuid.UUID, userID uuid.UUID, userRole models.Role, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to retrieve product for update")
	}
	if product == nil {
		return nil, apperror.ErrProductNotFound
	}

	// Authorization check: Only owner or admin can update
	if userRole != models.RoleAdmin && product.UserID != userID {
		return nil, apperror.ErrForbiddenProduct
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
		return nil, apperror.NewInternal(err, "failed to update product")
	}

	res := dto.ToProductResponse(product)
	return &res, nil
}

func (s *productService) Delete(id uuid.UUID, userID uuid.UUID, userRole models.Role) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return apperror.NewInternal(err, "failed to retrieve product for deletion")
	}
	if product == nil {
		return apperror.ErrProductNotFound
	}

	// Authorization check: Only owner or admin can delete
	if userRole != models.RoleAdmin && product.UserID != userID {
		return apperror.ErrForbiddenProduct
	}

	if err := s.productRepo.Delete(id); err != nil {
		return apperror.NewInternal(err, "failed to delete product")
	}

	return nil
}
