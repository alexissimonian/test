package auth

import (
	"fmt"
	"net/http"
	"strings"
)


func GetApiKey(header http.Header) (string, error) {
	tokenBearer := header.Get("Authorization")
	if len(tokenBearer) == 0 {
		return "", fmt.Errorf("No Authorization header")
	}

	token := strings.TrimPrefix(tokenBearer, "ApiKey ")
	return token, nil
}

