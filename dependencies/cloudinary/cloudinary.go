package cloudinary

// Import resty into your code and refer it as `resty`.
import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/go-resty/resty/v2"
	"net/http"

	"log"
	"time"
)

type Config struct {
	ApiKey    string
	ApiSecret string
	CloudName string
	BaseUrl   string
}

type uploadResponse struct {
	AssetId   string `json:"asset_id"`
	PublicId  string `json:"public_id"`
	SecureUrl string `json:"secure_url"`
	Url       string `json:"url"`
}

func (c Config) UploadFile(file []byte, filename string) (string, error) {

	client := resty.New()

	now := time.Now()
	signature, _ := c.getSignature(filename, now)

	endpoint := fmt.Sprintf("%s/%s/image/upload", c.BaseUrl, c.CloudName)
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetFileReader("file", "filename", bytes.NewReader(file)).
		SetFormData(map[string]string{
			"public_id": filename,
			"signature": signature,
			"api_key":   c.ApiKey,
			"timestamp": fmt.Sprintf("%d", now.Unix()),
		}).
		SetResult(uploadResponse{}).
		SetBasicAuth(c.ApiKey, c.ApiSecret).
		Post(endpoint)

	if err != nil {
		log.Println("Unable to upload image")
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		err = errors.New(fmt.Sprintf("Unable to upload image: status code %d", resp.StatusCode()))
		log.Println(err)
		return "", err
	}

	response := resp.Result().(*uploadResponse)

	return response.SecureUrl, nil
}

func (c Config) getSignature(filename string, now time.Time) (string, error) {
	themSigned, err := api.SignParameters(
		map[string][]string{
			"public_id": {filename},
			"timestamp": {fmt.Sprintf("%d", now.Unix())},
		},
		c.ApiSecret,
	)

	return themSigned, err
}
