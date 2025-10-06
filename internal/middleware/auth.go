package middleware

import (
	"encoding/json"
	"net/http"
	"practice_2/internal/i18n"
	"practice_2/internal/logger"
)

func AuthMiddleware(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := i18n.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
			msg := i18n.Get(lang)

			logger.Info("%s %s", r.Method, r.URL.Path)

			apiKey := r.Header.Get("X-API-Key")
			if apiKey != expectedKey {
				logger.Warn("Unauthorized access attempt from %s", r.RemoteAddr)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": msg.Unauthorized})
				return
			}

			logger.Debug("Request authenticated successfully")
			next.ServeHTTP(w, r)
		})
	}
}
