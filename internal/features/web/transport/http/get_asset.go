package web_transport_http

import (
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

func (h *handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	assetPath := r.PathValue("filepath")
	asset, err := h.webService.GetAsset(assetPath)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get static asset")
		return
	}

	respHandler.StaticAssetResponse(assetPath, asset)
}
