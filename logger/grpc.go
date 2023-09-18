package logger

import (
	"os"

	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

func init() {
	logLevel := getGRPCLogLevel()
	logger, err := newLogger(os.Stderr, logLevel)
	if err != nil {
		panic(err)
	}
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger.Logger))
}

func getGRPCLogLevel() zapcore.Level {
	// GRPC_GO_LOG_SEVERITY_LEVEL
	// see https://github.com/grpc/grpc-go/blob/v1.58.1/grpclog/loggerv2.go#L146-L154
	severity := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch severity {
	case "", "ERROR", "error": // If env is unset, set level to ERROR.
		return zapcore.ErrorLevel
	case "WARNING", "warning":
		return zapcore.WarnLevel
	case "INFO", "info":
		return zapcore.InfoLevel
	default:
		return zapcore.ErrorLevel
	}
}
