package services

import (
	"demerzel-events/internal/models"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func UpdateGroupService(tx *gorm.DB, req models.UpdateGroupRequest, id string) (int, models.Group, error) {
	group := models.Group{
		ID: id,
	}

	err := group.GetGroupByID(tx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return http.StatusNotFound, group, fmt.Errorf("group with the specified id does not exist")
		}
		return http.StatusBadRequest, group, err
	}

	// check if a `name` was passed in the request body
	if req.Name != "" {
		group.Name = req.Name
	}

	// update group id
	err = group.UpdateGroupByID(tx)
	if err != nil {
		return http.StatusInternalServerError, group, err
	}

	return http.StatusOK, group, nil
}
