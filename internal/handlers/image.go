package handlers

import (
	"bytes"
	"demerzel-events/dependencies/cloudinary"
	"demerzel-events/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func UploadFileHandler(c *gin.Context) {
	uploadedFile, _ := c.FormFile("file")
	if uploadedFile == nil {
		response.Error(c, http.StatusBadRequest, "No files specified")
		return
	}

	file, err := uploadedFile.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "Unable to read upload:"+err.Error())
		return
	}

	uploader := cloudinary.Config{
		ApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		ApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		BaseUrl:   os.Getenv("CLOUDINARY_BASE_URL"),
	}
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), uploadedFile.Filename)

	url, err := uploader.UploadFile(buf.Bytes(), filename)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Unable to upload file:"+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "File uploaded", map[string]string{"url": url})
	return

}
