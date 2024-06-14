package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (apiCfg *ApiConfig) cloudinarUploader(r *http.Request, fileName string) (string, error) {
	// Upload a file
	uploadResult, err := apiCfg.Cld.Upload.Upload(r.Context(), fmt.Sprintf("uploads/%v", fileName), uploader.UploadParams{})
	if err != nil {
		fmt.Println("cloudinary error", err)
		return "", errors.New("Error uploading to cloudinary")
	}

	return uploadResult.SecureURL, nil
}
