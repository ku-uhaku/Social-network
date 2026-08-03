package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// SaveUploadedImage saves a multipart uploaded image to disk and returns the stored filename.
func SaveUploadedImage(file multipart.File, header *multipart.FileHeader, mediaDir string) (string, error) {
	if file == nil || header == nil || header.Size == 0 {
		return "", fmt.Errorf("uploaded file is empty")
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read uploaded file: %w", err)
	}

	hash := sha256.Sum256(fileBytes)
	filename := hex.EncodeToString(hash[:])

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = http.DetectContentType(fileBytes)
		}
		if guesses, _ := mime.ExtensionsByType(contentType); len(guesses) > 0 {
			ext = strings.ToLower(guesses[0])
		}
	}
	if ext == "" {
		ext = ".bin"
	}

	filename += ext
	if err := os.MkdirAll(mediaDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to prepare media directory: %w", err)
	}

	path := filepath.Join(mediaDir, filename)
	if err := os.WriteFile(path, fileBytes, 0o644); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}

	return filename, nil
}
