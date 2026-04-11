package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/rod1kutzyy/task-manager-app/internal/core/config"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	core_pgx_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool/pgx"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/middleware"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/rod1kutzyy/task-manager-app/internal/features/statistics/repository/postgres"
	statistics_service "github.com/rod1kutzyy/task-manager-app/internal/features/statistics/service"
	statistics_transport_http "github.com/rod1kutzyy/task-manager-app/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/repository/postgres"
	tasks_service "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/service"
	tasks_transport_http "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/rod1kutzyy/task-manager-app/internal/features/users/repository/postgres"
	users_service "github.com/rod1kutzyy/task-manager-app/internal/features/users/service"
	users_transport_http "github.com/rod1kutzyy/task-manager-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()

	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		log.Fatalf("failed to init app logger: %v", err)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewRepository(pool)
	usersService := users_service.NewService(usersRepository)
	usersTransportHTTP := users_transport_http.NewHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewRepository(pool)
	tasksService := tasks_service.NewService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewRepository(pool)
	statisticsService := statistics_service.NewService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewHandler(statisticsService)

	logger.Debug("initializing HTTP server")
	httpServer := server.NewHTTPServer(
		server.NewConfigMust(),
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.Recovery(),
	)

	apiVersionRouter := server.NewAPIVersionRouter(server.ApiVersion1)

	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
