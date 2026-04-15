package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	w      http.ResponseWriter
	logger *core_logger.Logger
}

func NewHTTPResponseHandler(w http.ResponseWriter, logger *core_logger.Logger) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		w:      w,
		logger: logger,
	}
}

func (h *HTTPResponseHandler) JSONResponse(body any, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(body); err != nil {
		h.logger.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) HTMLResponse(html []byte) {
	h.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.w.WriteHeader(http.StatusOK)

	if _, err := h.w.Write(html); err != nil {
		h.logger.Error("write HTML HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) StaticAssetResponse(assetPath string, asset []byte) {
	contentType := mime.TypeByExtension(filepath.Ext(assetPath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h.w.Header().Set("Content-Type", contentType)
	h.w.WriteHeader(http.StatusOK)

	if _, err := h.w.Write(asset); err != nil {
		h.logger.Error("write static asset HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.logger.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.logger.Debug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.logger.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.logger.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.logger.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	resp := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.JSONResponse(resp, statusCode)
}
