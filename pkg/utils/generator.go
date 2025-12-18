package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// Generate a random short code
func GenerateShortCode (length int) (string, error){
	bytes := make([]byte, length)
	if _, err:= rand.Read(bytes); err != nil {
		return "", err
	}
	// URL-safe base64 encoding and remove padding
	encoded := base64.URLEncoding.EncodeToString(bytes)
	encoded = strings.TrimRight(encoded, "=")
	if len(encoded) > length {
		encoded = encoded[:length]
	}
	return encoded, nil
}

// Check if the URL is valid
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func ValidateShortCode(code string) bool {
	if len(code) < 3 || len(code) > 10{
		return false
	}
	
	for _, c := range code {
		if ! ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}