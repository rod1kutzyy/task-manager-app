package web_service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
)

func (s *service) GetAsset(assetPath string) ([]byte, error) {
	cleanPath := filepath.Clean(assetPath)

	if cleanPath == "." || strings.HasPrefix(cleanPath, "..") || filepath.IsAbs(cleanPath) {
		return nil, fmt.Errorf("invalid asset path %q: %w", assetPath, core_errors.ErrInvalidArgument)
	}

	filePath := filepath.Join(os.Getenv("PROJECT_ROOT"), "public", "assets", cleanPath)

	asset, err := s.webRepository.GetFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("get asset from repository: %w", err)
	}

	return asset, nil
}
