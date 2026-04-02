package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/middleware"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
	users_transport_http "github.com/rod1kutzyy/task-manager-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		log.Fatal("failed to init app logger: %w", err)
	}
	defer logger.Close()

	logger.Debug("Starting app!")

	usersTransportHTTP := users_transport_http.NewHandler(nil)

	usersRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := server.NewAPIVersionRouter(server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := server.NewHTTPServer(
		server.NewConfigMust(),
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(),
		middleware.Trace(),
	)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
