package web_transport_http

import (
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

func (h *handler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	html, err := h.webService.GetMainPage()
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get index.html for main page")
		return
	}

	respHandler.HTMLResponse(html)
}
