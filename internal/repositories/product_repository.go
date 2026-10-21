package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/internal/models"
	"github.com/Maulighadge/GolangTemplate/pkg/utils"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uuid.UUID) (*models.Product, error)
	FindBySKU(sku string) (*models.Product, error)
	FindAll(search, category string, pagination utils.Pagination) ([]models.Product, int64, error)
	Update(product *models.Product) error
	Delete(id uuid.UUID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("User").First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "sku = ?", sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(search, category string, pagination utils.Pagination) ([]models.Product, int64, error) {
	var products []models.Product
	var totalRows int64

	query := r.db.Model(&models.Product{})

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR sku LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("User").
		Order("created_at DESC").
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, totalRows, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}
