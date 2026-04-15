package web_transport_http

import (
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
)

type handler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewHandler(webService WebService) *handler {
	return &handler{
		webService: webService,
	}
}

func (h *handler) Routes() []server.Route {
	return []server.Route{
		{
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
