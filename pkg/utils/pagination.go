package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination meta information
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationFromContext extracts page and limit from query parameters with defaults
func GetPaginationFromContext(c *gin.Context) (page int, limit int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100 // Cap limit at 100
	}

	return page, limit
}

// CalculatePagination creates Pagination struct with total calculations
func CalculatePagination(page, limit int, totalRows int64) Pagination {
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}
}

// Offset calculates the SQL query offset
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
