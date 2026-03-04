package config

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary(cfg *Config) (*cloudinary.Cloudinary, error) {
	url := fmt.Sprintf(
		"cloudinary://%s:%s@%s",
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
		cfg.CloudinaryCloudName,
	)

	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, err
	}

	// Test koneksi
	if _, err := cld.Admin.Ping(context.Background()); err != nil {
		return nil, err
	}

	return cld, nil
}