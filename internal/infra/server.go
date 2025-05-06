package infra

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/clubhub.ai1/gommon/logger"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"go.uber.org/fx"
)

var (
	readTimeout  = time.Second * 30
	writeTimeout = time.Second * 30
)

func NewHTTPServer() *http.Server {
	return router.GetServerConfiguration(config.Config().Server.Port,
		router.WithReadTimeout(readTimeout),
		router.WithWriteTimeout(writeTimeout),
	)
}

func NewEchoServer(lc fx.Lifecycle,
	logger logger.LoggerInterface,
	httpServer *http.Server,
) *echo.Echo {
	server := router.New(router.Config{
		HealthPath:  config.Config().Server.BasePath + "/health",
		ServiceName: config.Config().App.ServiceName,
		Logger:      logger,
		NoTraceID:   false,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.StartServer(httpServer); err != nil {
					logger.Fatalf("Error starting Server: ", err)
				}

				logger.Infof("Server is listening on PORT: %d", config.Config().Server.Port)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}
