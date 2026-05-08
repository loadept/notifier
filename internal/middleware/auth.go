package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

func Auth(secret string) Middleware {
	secretBytes := []byte(secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				http.Error(w, "unsupported type", http.StatusUnsupportedMediaType)
				return
			}

			signature := r.Header.Get("X-Signature")
			timestamp := r.Header.Get("X-Timestamp")
			if signature == "" || timestamp == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			cleanPath := path.Clean(r.URL.Path)
			message := fmt.Sprintf("%s\n%s\n%s\n%s", r.Method, cleanPath, timestamp, body)

			mac := hmac.New(sha256.New, secretBytes)
			mac.Write([]byte(message))
			expected := hex.EncodeToString(mac.Sum(nil))
			if subtle.ConstantTimeCompare([]byte(signature), []byte(expected)) == 0 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
