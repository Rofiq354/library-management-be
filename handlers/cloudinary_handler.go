package handlers

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (h *AuthHandler) uploadToCloudinary(file *multipart.FileHeader) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	result, err := h.Cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{
		Folder: "library/books",
	})
	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}