package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Upload image helper
func (h *AuthHandler) uploadToCloudinary(file *multipart.FileHeader, folder string) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	result, err := h.Cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}


// Upload PDF helper - dengan resource_type: "raw"
func (h *AuthHandler) uploadPDFToCloudinary(file *multipart.FileHeader) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	// Validate PDF
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		return "", "", fmt.Errorf("file harus berformat PDF")
	}

	result, err := h.Cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{
		Folder:       "library/books/pdfs",
		ResourceType: "raw", // ← Penting untuk PDF
	})
	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}