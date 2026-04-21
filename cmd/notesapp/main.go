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
	core_goredis_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/redis/pool/goredis"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/middleware"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/rod1kutzyy/task-manager-app/internal/features/statistics/repository/postgres"
	statistics_service "github.com/rod1kutzyy/task-manager-app/internal/features/statistics/service"
	statistics_transport_http "github.com/rod1kutzyy/task-manager-app/internal/features/statistics/transport/http"
	tasks_service "github.com/rod1kutzyy/task-manager-app/internal/features/tasks"
	tasks_adapters_in_transport_http "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/adapters/in/transport/http"
	tasks_adapters_out_repository_cached "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/adapters/out/repository/cached"
	tasks_adapters_out_repository_postgres "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/adapters/out/repository/postgres"
	users_postgres_repository "github.com/rod1kutzyy/task-manager-app/internal/features/users/repository/postgres"
	users_service "github.com/rod1kutzyy/task-manager-app/internal/features/users/service"
	users_transport_http "github.com/rod1kutzyy/task-manager-app/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/rod1kutzyy/task-manager-app/docs"
)

// @title Notes app API
// @version 1.0
// @description REST API for managing users, tasks, and task statistics.
// @BasePath /api/v1
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
	postgresPool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer postgresPool.Close()

	logger.Debug("initializing redis connection pool")
	redisPool, err := core_goredis_pool.NewPool(ctx, core_goredis_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init redis connection pool", zap.Error(err))
	}
	defer redisPool.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewRepository(postgresPool)
	usersService := users_service.NewService(usersRepository)
	usersTransportHTTP := users_transport_http.NewHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_adapters_out_repository_cached.NewCachedRepository(
		redisPool,
		tasks_adapters_out_repository_postgres.NewRepository(postgresPool),
	)
	tasksService := tasks_service.NewService(tasksRepository)
	tasksTransportHTTP := tasks_adapters_in_transport_http.NewHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewRepository(postgresPool)
	statisticsService := statistics_service.NewService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewHandler(statisticsService)

	logger.Debug("initializing HTTP server")
	httpConfig := server.NewConfigMust()
	httpServer := server.NewHTTPServer(
		httpConfig,
		logger,
		middleware.CORS(httpConfig.AllowedOrigins),
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
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
