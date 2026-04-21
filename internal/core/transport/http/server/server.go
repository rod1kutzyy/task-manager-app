package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/docs"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux         *http.ServeMux
	config      Config
	logger      *core_logger.Logger
	middlewares []middleware.Middleware
}

func NewHTTPServer(config Config, logger *core_logger.Logger, middlewares ...middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      config,
		logger:      logger,
		middlewares: middlewares,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddlewares()))
	}
}

func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		s.mux.Handle(pattern, route.WithMiddlewares())
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		swaggerInfo := *docs.SwaggerInfo
		swaggerInfo.Host = r.Host

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(swaggerInfo.ReadDoc()))
	})
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(s.mux, s.middlewares...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		s.logger.Warn(
			"start HTTP server",
			zap.String("addr", s.config.Addr),
		)

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.logger.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.logger.Warn("HTTP server stopped")
	}

	return nil
}
