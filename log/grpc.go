package log

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"google.golang.org/grpc/grpclog"
)

var _ grpclog.LoggerV2 = (*grpcLogger)(nil)

type grpcLogger struct {
	h        slog.Handler
	minLevel slog.Level
}

func NewGRPCLogger(h slog.Handler) *grpcLogger {
	return &grpcLogger{h: h, minLevel: getGRPCLogLevel()}
}

func (l *grpcLogger) Info(args ...any) {
	l.log(slog.LevelInfo, fmt.Sprint(args...))
}

func (l *grpcLogger) Infoln(args ...any) {
	l.log(slog.LevelInfo, fmt.Sprintln(args...))
}

func (l *grpcLogger) Infof(format string, args ...any) {
	l.log(slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Warning(args ...any) {
	l.log(slog.LevelWarn, fmt.Sprint(args...))
}

func (l *grpcLogger) Warningln(args ...any) {
	l.log(slog.LevelWarn, fmt.Sprintln(args...))
}

func (l *grpcLogger) Warningf(format string, args ...any) {
	l.log(slog.LevelWarn, fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Error(args ...any) {
	l.log(slog.LevelError, fmt.Sprint(args...))
}

func (l *grpcLogger) Errorln(args ...any) {
	l.log(slog.LevelError, fmt.Sprintln(args...))
}

func (l *grpcLogger) Errorf(format string, args ...any) {
	l.log(slog.LevelError, fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Fatal(args ...any) {
	l.log(slog.LevelError+1, fmt.Sprint(args...))
}

func (l *grpcLogger) Fatalln(args ...any) {
	l.log(slog.LevelError+1, fmt.Sprintln(args...))
}

func (l *grpcLogger) Fatalf(format string, args ...any) {
	l.log(slog.LevelError+1, fmt.Sprintf(format, args...))
}

func (l *grpcLogger) V(level int) bool {
	return l.h.Enabled(context.Background(), slog.Level(int(slog.LevelDebug)-level))
}

func (l *grpcLogger) log(level slog.Level, msg string) {
	if l.minLevel > level {
		return
	}
	//nolint: dogsled
	pc, _, _, _ := runtime.Caller(2)
	r := slog.NewRecord(time.Now(), level, msg, pc)
	if err := l.h.Handle(context.Background(), r); err != nil {
		log.Printf("log: failed to handle record in grpcLogger: %v", err)
	}
}

func getGRPCLogLevel() slog.Level {
	// GRPC_GO_LOG_SEVERITY_LEVEL
	// see https://github.com/grpc/grpc-go/blob/v1.58.1/grpclog/loggerv2.go#L146-L154
	severity := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch severity {
	// If env is unset, set level to ERROR.
	case "", "ERROR", "error":
		return slog.LevelError
	case "WARNING", "warning":
		return slog.LevelWarn
	case "INFO", "info":
		return slog.LevelInfo
	default:
		return slog.LevelError
	}
}
