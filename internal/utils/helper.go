package utils

import "github.com/truongbo17/go-gin-boilerplate/config"

func CheckLanguage(target string) bool {
	for _, code := range config.LanguageAvailable {
		if code == target {
			return true
		}
	}
	return false
}
