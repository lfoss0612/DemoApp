package server

import (
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

