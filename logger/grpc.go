package logger

import (
	"os"

	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

func init() {
	logLevel := getGRPCLogLevel()
	logger := newLogger(os.Stderr, logLevel)
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger.Logger))
}

func getGRPCLogLevel() zapcore.Level {
	// GRPC_GO_LOG_SEVERITY_LEVEL
	// see https://github.com/grpc/grpc-go/blob/v1.58.1/grpclog/loggerv2.go#L146-L154
	severity := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch severity {
	// If env is unset, set level to ERROR.
	//nolint:goconst
	case "", "ERROR", "error":
		return zapcore.ErrorLevel
	case "WARNING", "warning":
		return zapcore.WarnLevel
	//nolint:goconst
	case "INFO", "info":
		return zapcore.InfoLevel
	default:
		return zapcore.ErrorLevel
	}
}
