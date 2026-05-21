package cloudinary

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service struct {
	cld *cloudinary.Cloudinary
}

func New(cloudName, apiKey, apiSecret string) (*Service, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("cloudinary init: %w", err)
	}
	return &Service{cld: cld}, nil
}

func (s *Service) Upload(ctx context.Context, file interface{}, filename string) (string, error) {
	uniqueFilename := true
	overwrite := false
	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       filename,
		UniqueFilename: &uniqueFilename,
		Overwrite:      &overwrite,
		Folder:         "rizzra",
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary upload: %w", err)
	}
	return result.SecureURL, nil
}
