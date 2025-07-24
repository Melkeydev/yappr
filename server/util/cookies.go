package util

import (
	"net/http"
)

// SetSecureCookie sets a cookie with appropriate security settings based on environment
func SetSecureCookie(w http.ResponseWriter, name, value string, maxAge int) {
	env := GetEnv("ENVIRONMENT", "dev")
	
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	}

	if env != "prod" {
		// Development environment
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	} else {
		// Production environment
		cookie.Domain = ".yappr.chat"
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
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

	if env != "prod" {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	} else {
		cookie.Domain = ".yappr.chat"
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, cookie)
}