package utils

import (
	"os"
	"strings"
)

func DetectLanguage() string {
	lang := os.Getenv("LC_ALL")
	if lang == "" {
		lang = os.Getenv("LANG")
	}
	if lang == "" {
		return "en"
	}
	
	parts := strings.Split(lang, ".")
	if len(parts) > 0 {
		lang = parts[0]
	}
	parts = strings.Split(lang, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return "en"
}



