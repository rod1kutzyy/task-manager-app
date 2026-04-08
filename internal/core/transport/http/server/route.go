package server

import (
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/middleware"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []middleware.Middleware
}

func (r *Route) WithMiddlewares() http.Handler {
	return middleware.ChainMiddleware(r.Handler, r.Middlewares...)
}
