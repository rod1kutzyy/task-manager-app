package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	"github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	logger *logger.Logger
	w      http.ResponseWriter
}

func NewHTTPResponseHandler(logger *logger.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		logger: logger,
		w:      w,
	}
}

func (h *HTTPResponseHandler) JSONResponse(body any, statusCode int) {
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(body); err != nil {
		h.logger.Error("write HTTP response", zap.Error(err))
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
	resp := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(resp, statusCode)
}
