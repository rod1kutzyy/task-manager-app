package response

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.logger.Error(msg, zap.Error(err))
	h.w.WriteHeader(statusCode)

	resp := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.w).Encode(resp); err != nil {
		h.logger.Error("write HTTP response", zap.Error(err))
	}
}
