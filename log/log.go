package log

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
)

const ErrorKey = "error"

func SetOption(opts ...Option) error {
	return SetOptionWithWriter(os.Stderr, opts...)
}

func SetOptionWithWriter(w io.Writer, opts ...Option) error {
	logger := NewCloudLogging(w, opts...)
	slog.SetDefault(logger)
	return nil
}

func Err(err error) slog.Attr {
	return slog.Any(ErrorKey, err)
}

func Source(r slog.Record) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	return &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}

func StackTrace(r slog.Record) string {
	stack := debug.Stack()
	source := Source(r)
	return chopStack(stack, source.Function)
}

// chopStack trims a stack trace so that the function which panics or calls Error is first.
// original: https://github.com/googleapis/google-cloud-go/blob/errorreporting/v0.3.0/errorreporting/errors.go#L211-L234
func chopStack(s []byte, target string) string {
	headerLine := bytes.IndexByte(s, '\n')
	if headerLine == -1 {
		return string(s)
	}
	stack := s[headerLine:]
	targetLine := bytes.Index(stack, []byte(target))
	if targetLine == -1 {
		return string(s)
	}
	stack = stack[targetLine:]
	return string(s[:headerLine+1]) + string(stack)
}

func LevelFromEnv(env string) slog.Level {
	logLevel := os.Getenv(env)
	switch logLevel {
	case "DEBUG", "debug":
		return slog.LevelDebug
	case "INFO", "info":
		return slog.LevelInfo
	case "WARN", "warn":
		return slog.LevelWarn
	case "ERROR", "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
