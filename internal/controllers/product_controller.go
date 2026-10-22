package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/internal/dto"
	"github.com/Maulighadge/GolangTemplate/internal/middlewares"
	"github.com/Maulighadge/GolangTemplate/internal/services"
	"github.com/Maulighadge/GolangTemplate/pkg/apperror"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
	"github.com/Maulighadge/GolangTemplate/pkg/utils"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

// Create creates a new product
func (ctl *ProductController) Create(c *gin.Context) {
	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized access")
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	product, err := ctl.productService.Create(userID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, "Product created successfully", product)
}

// GetByID retrieves a product by its UUID
func (ctl *ProductController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.HandleError(c, apperror.ErrInvalidID)
		return
	}

	product, err := ctl.productService.GetByID(id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, "Product retrieved successfully", product)
}

// List retrieves paginated products with optional search and category filters
func (ctl *ProductController) List(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	page, limit := utils.GetPaginationFromContext(c)

	products, meta, err := ctl.productService.List(search, category, page, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Paginated(c, "Products retrieved successfully", products, meta)
}

// Update updates an existing product
func (ctl *ProductController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.HandleError(c, apperror.ErrInvalidID)
		return
	}

	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized access")
		return
	}
	role, _ := middlewares.GetAuthUserRole(c)

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	product, err := ctl.productService.Update(id, userID, role, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, "Product updated successfully", product)
}

// Delete removes a product
func (ctl *ProductController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.HandleError(c, apperror.ErrInvalidID)
		return
	}

	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized access")
		return
	}
	role, _ := middlewares.GetAuthUserRole(c)

	if err := ctl.productService.Delete(id, userID, role); err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, "Product deleted successfully", nil)
}
