package handlers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"demerzel-events/dependencies/cloudinary"
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"

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
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid request body format: %s", err.Error()),
		)
		return
	}

	code, data, err := services.UpdateGroupService(db.DB, req, id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Group updated successfully", data)
}

func ListGroups(c *gin.Context) {
	name := c.DefaultQuery("name", "")

	f := services.Filter{
		Search: struct{ Name string }{
			Name: name,
		},
	}

	groups, err := services.ListGroups(f)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error: failed to fetch groups")
		return
	}

	var message string

	if len(groups) == 0 {
		message = "No groups"
	}

	if len(groups) > 0 {
		message = "Groups retrieved successfully"
	}

	response.Success(c, http.StatusOK, message, groups)
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

	userGroups, code, err := services.GetGroupsByUserId(user.Id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Fetched all user groups", userGroups)
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
	err := services.DeleteGroup(db.DB, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, fmt.Sprintf("group with id=%s deleted successfully", id), nil)

}

func GroupEventsById(c *gin.Context) {
	id := c.Param("id")

	group := models.Group{ID: id}
	result, err := group.GetGroupEvents(db.DB)

	if err != nil {
		response.Error(c, 500, "Can't process your request")
		return
	}

	response.Success(c, 200, "Group events retrieved", result)
}
