package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// GetParamInt64 extracts a query parameter from the URL string and parses it to int64
func GetParamInt64(r *http.Request, key string) (int64, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return 0, errors.New("missing parameter: " + key)
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid parameter format: " + key)
	}

	return val, nil
}

const maxImageUploadSize = 20 << 20 // 20 MB

// MediaDir is where uploaded images are stored and served from.
const MediaDir = "media"

var allowedImageTypes = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

func SaveUploadedImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	if file == nil || header == nil || header.Size == 0 {
		return "", fmt.Errorf("uploaded file is empty")
	}

	if header.Size > maxImageUploadSize {
		return "", fmt.Errorf("uploaded file exceeds maximum allowed size of 20MB")
	}

	fileBytes, err := io.ReadAll(io.LimitReader(file, maxImageUploadSize+1))
	if err != nil {
		return "", fmt.Errorf("failed to read uploaded file: %w", err)
	}

	if int64(len(fileBytes)) > maxImageUploadSize {
		return "", fmt.Errorf("uploaded file exceeds maximum allowed size of 20MB")
	}

	// check header instead of extension
	detectedType := http.DetectContentType(fileBytes)
	ext, ok := allowedImageTypes[detectedType]
	if !ok {
		return "", fmt.Errorf("unsupported image type %q: only PNG, JPEG, and GIF are allowed", detectedType)
	}

	hash := sha256.Sum256(fileBytes)
	filename := hex.EncodeToString(hash[:]) + ext

	if err := os.MkdirAll(MediaDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to prepare media directory: %w", err)
	}

	path := filepath.Join(MediaDir, filename)
	if err := os.WriteFile(path, fileBytes, 0o644); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}

	return filename, nil
}
