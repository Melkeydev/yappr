package util

import (
	"net/http"
)

func SetSecureCookie(w http.ResponseWriter, name, value string, maxAge int) {
	env := GetEnv("ENVIRONMENT", "dev")

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	}

	if env == "prod" {
		// Production environment
		cookie.Domain = ".yappr.chat"
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		// Development environment
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, cookie)
}

// ClearSecureCookie clears a cookie with appropriate security settings
func ClearSecureCookie(w http.ResponseWriter, name string) {
	env := GetEnv("ENVIRONMENT", "dev")

	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	if env == "prod" {
		cookie.Domain = ".yappr.chat"
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, cookie)
}

