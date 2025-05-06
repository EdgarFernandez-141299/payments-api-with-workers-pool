package logger

import (
	traceApm "gitlab.com/clubhub.ai1/go-libraries/observability/apm"
	logger2 "gitlab.com/clubhub.ai1/go-libraries/observability/logger"
	"gitlab.com/clubhub.ai1/gommon/logger"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/observability/adapters"
	"os"
)

func NewLoggerLegacy() logger.LoggerInterface {
	loggerSettings := logger.LoggerSettings{
		MicroserviceName: config.Config().App.ServiceName,
		ConsoleLog: logger.ConsoleLoggerSettings{
			DebugMode: config.Config().App.LoggerDebugMode,
		},
	}

	return logger.NewLogger(loggerSettings)
}

func NewLogger() adapters.Logger {
	serviceName := config.Config().App.ServiceName

	log := logger2.New(os.Stdout, logger2.LevelInfo, serviceName, nil)

	traceApm.NewTraceProvider(log)

	return log
}

func NewTraceProvider(log adapters.Logger) traceApm.Tracer {
	traceApm.NewTraceProvider(log)
	return traceApm.NewTracer("payments-api")
}
