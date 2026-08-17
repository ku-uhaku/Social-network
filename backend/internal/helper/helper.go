package helper

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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

const (
	maxImageUploadSize = 20 << 20 // 20 MB
	maxImageDimension  = 8000     // max width/height per side, guards against decompression bombs
)

// MediaDir is where uploaded images are stored and served from.
const MediaDir = "media"

// imageExtensions maps a decoded image format to its stored file extension.
var imageExtensions = map[string]string{
	"jpeg": ".jpg",
	"png":  ".png",
	"gif":  ".gif",
}


// avoid renaaaame file 
// the deteeeected of foormaat "jpeg", "png", or "gif"
func IsValidImage(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empttty  image data")
	}

	// DecodeConfig validates the header and reports the format without
	// alllocating a fullll pixel buffer
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("not a valid image: %w", err)
	}

	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > maxImageDimension || cfg.Height > maxImageDimension {
		return "", fmt.Errorf("image dimensions %dx%d exceed the allowed limit", cfg.Width, cfg.Height)
	}

	if _, ok := imageExtensions[format]; !ok {
		return "", fmt.Errorf("unsupported image format %q: only PNG, JPEG, and GIF are allowed", format)
	}

	// Fully decode to reject files with a valid header but a corrupt body.
	if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("corrupt image data: %w", err)
	}

	return format, nil
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

	// Validate the bytes are a real image and pick the matching extension.
	format, err := IsValidImage(fileBytes)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileBytes)
	filename := hex.EncodeToString(hash[:]) + imageExtensions[format]

	if err := os.MkdirAll(MediaDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to prepare media directory: %w", err)
	}

	path := filepath.Join(MediaDir, filename)
	if err := os.WriteFile(path, fileBytes, 0o644); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}

	return filename, nil
}
