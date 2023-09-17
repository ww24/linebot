package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

func init() {
	logLevel := getGRPCLogLevel()
	logger, err := new(os.Stderr)
	if err != nil {
		panic(err)
	}
	l := logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		newCore, err := zapcore.NewIncreaseLevelCore(core, logLevel)
		if err != nil {
			logger.Warn("logger: failed to increase log level",
				zap.Error(err),
				zap.String("level", logLevel.String()),
			)
			return core
		}
		return newCore
	}))
	grpclog.SetLoggerV2(zapgrpc.NewLogger(l))
}

func getGRPCLogLevel() zapcore.Level {
	// GRPC_GO_LOG_SEVERITY_LEVEL
	// see https://github.com/grpc/grpc-go/blob/v1.58.1/grpclog/loggerv2.go#L146-L154
	severity := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch severity {
	case "", "ERROR", "error": // If env is unset, set level to ERROR.
		return zap.ErrorLevel
	case "WARNING", "warning":
		return zap.WarnLevel
	case "INFO", "info":
		return zap.InfoLevel
	default:
		return zap.ErrorLevel
	}
}
