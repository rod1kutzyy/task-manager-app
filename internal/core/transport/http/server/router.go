package server

import (
	"fmt"
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion  ApiVersion
	middlewares []middleware.Middleware
}

func NewAPIVersionRouter(apiVersion ApiVersion, middlewares ...middleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:    http.NewServeMux(),
		apiVersion:  apiVersion,
		middlewares: middlewares,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddlewares())
	}
}

func (r *APIVersionRouter) WithMiddlewares() http.Handler {
	return middleware.ChainMiddleware(r, r.middlewares...)
}
