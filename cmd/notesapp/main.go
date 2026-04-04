package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/middleware"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
	users_postgres_repository "github.com/rod1kutzyy/task-manager-app/internal/features/users/repository/postgres"
	users_service "github.com/rod1kutzyy/task-manager-app/internal/features/users/service"
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

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewRepository(pool)
	usersService := users_service.NewService(usersRepository)
	usersTransportHTTP := users_transport_http.NewHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := server.NewHTTPServer(
		server.NewConfigMust(),
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(),
		middleware.Trace(),
	)

	apiVersionRouter := server.NewAPIVersionRouter(server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
