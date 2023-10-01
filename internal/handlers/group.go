package handlers

import (
	"bytes"
	"demerzel-events/dependencies/cloudinary"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/helpers"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {
	type jsonData struct {
		Name string `json:"name" binding:"required"`
	}

	var requestBody struct {
		File      *multipart.FileHeader `form:"file"`
		GroupData jsonData              `form:"jsonData"`
	}

	if err := ctx.ShouldBind(&requestBody); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid request body format: %s", err.Error()),
		)
		return
	}

	var photUrl string = ""
	fileToUpload := requestBody.File

	if fileToUpload != nil {

		image, err := fileToUpload.Open()
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error())
		}

		buffer := new(bytes.Buffer)
		_, err = buffer.ReadFrom(image)

		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "Cannot process file:"+err.Error())
			return
		}

		transport := cloudinary.Config{
			ApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
			ApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
			CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
			BaseUrl:   os.Getenv("CLOUDINARY_BASE_URL"),
		}

		imageUrl, err := transport.UploadFile(buffer.Bytes(), fileToUpload.Filename)

		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "Could not upload to file to bucket:"+err.Error())
			return
		}

		photUrl = imageUrl
	}

	var newGroup models.Group
	newGroup.Name = requestBody.GroupData.Name
	newGroup.Image = photUrl

	services.CreateGroup(&newGroup)

	response.Success(
		ctx,
		http.StatusCreated,
		"Group created successfully",
		newGroup,
	)
}

func UpdateGroup(c *gin.Context) {
	req := models.UpdateGroupRequest{}
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid request body format: %s", err.Error()),
		)
		return
	}

	data, code, err := services.UpdateGroupById(id, &req)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Group updated successfully", data)
}

func ListGroups(c *gin.Context) {
	name := c.DefaultQuery("name", "")

	// Extract query parameters for pagination
	limit, offset, err := helpers.GetLimitAndOffset(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	groups, totalGroups, err := services.ListGroups(name, *limit, *offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error: failed to fetch groups")
		return
	}

	response.Success(c, http.StatusOK, "Groups retrieved successfully", map[string]interface{}{
		"groups":       groups,
		"total_groups": totalGroups,
	})
}

// GetUserGroups returns all group this user belongs to
func GetUserGroups(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	// Extract query parameters for pagination
	limit, offset, err := helpers.GetLimitAndOffset(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userGroups, totalUserGroups, err := services.GetGroupsByUserId(user.Id, *limit, *offset)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Fetched all user groups", map[string]interface{}{
		"user_groups":       userGroups,
		"total_user_groups": totalUserGroups,
	})
}

func GetGroupById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Group ID is required")
		return
	}

	group, err := services.GetGroupById(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Group retrieved successfully", group)
}

func DeleteGroup(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Please provide a valid group id")
		return
	}
	code, err := services.DeleteGroup(id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, fmt.Sprintf("group with id=%s deleted successfully", id), nil)

}

func GroupEventsById(c *gin.Context) {
	id := c.Param("id")

	result, code, err := services.GetEventByID(id)

	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Group events retrieved", result)
}
