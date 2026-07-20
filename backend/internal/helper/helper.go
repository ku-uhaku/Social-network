package helper

import (
	"errors"
	"net/http"
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
