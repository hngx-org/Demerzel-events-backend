package helpers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Extract and validate query parameters for pagination
func GetLimitAndOffset(c *gin.Context) (*int, *int, error) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Convert page and perPage parameters to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, nil, errors.New("Invalid page parameter")
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, nil, errors.New("Invalid limit parameter")
	}

	offset := (pageInt - 1) * limitInt
	return &limitInt, &offset, nil
}
