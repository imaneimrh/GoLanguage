package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func ExtractPathParamInt(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		return 0, errors.New("Url parameter not provided")
	}
	ID, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, errors.New("Invalid Url parameter")
	}
	return ID, nil
}
