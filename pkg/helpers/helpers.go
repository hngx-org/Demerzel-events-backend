package helpers

import (
	"errors"
	"strconv"
	"time"

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

func IsValidDate(date string) bool {
	layout := "2006-01-02"
	_, err := time.Parse(layout, date)

	return err == nil
}

func FormatDateTimeStr(dateStr string, timeStr string) (string, error) {
	dateTimeStr := dateStr + " " + timeStr

	dateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return "", err
	}

	return dateTime.Format("20060102T150405"), nil
}
