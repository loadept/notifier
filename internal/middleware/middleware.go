// Package middleware proporciona middlewares de diversas responsabilidades
// como auth o log
package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler
