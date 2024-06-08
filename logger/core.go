package logger

import (
	"bytes"
	"runtime/debug"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// core implements zapcore.Core.
type core struct {
	zapcore.Core
}

func newCore(c zapcore.Core) *core {
	return &core{
		Core: c,
	}
}

func (c *core) With(fields []zap.Field) zapcore.Core {
	return newCore(c.Core.With(fields))
}

//nolint:gocritic
func (c *core) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}

//nolint:gocritic
func (c *core) Write(e zapcore.Entry, fields []zapcore.Field) error {
	if zapcore.ErrorLevel.Enabled(e.Level) {
		if report := newErrorReport(e.Caller); report != nil {
			fields = append(fields,
				zap.Object("context", report),
				// see: https://cloud.google.com/error-reporting/docs/formatting-error-messages#log-error
				// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/projects.events/report#ReportedErrorEvent
				zap.String("stack_trace", e.Message+"\n"+chopStack(debug.Stack())),
			)
		}
	}
	if loc := newSourceLocation(e.Caller); loc != nil {
		fields = append(fields, zap.Object("logging.googleapis.com/sourceLocation", loc))
	}

	//nolint: wrapcheck
	return c.Core.Write(e, fields)
}

// chopStack trims a stack trace so that the function which panics or calls Error is first.
// original: https://github.com/googleapis/google-cloud-go/blob/errorreporting/v0.3.0/errorreporting/errors.go#L211-L234
func chopStack(s []byte) string {
	targets := []string{
		"go.uber.org/zap.(*Logger).Error",
		"go.uber.org/zap.(*Logger).DPanic",
		"go.uber.org/zap.(*Logger).Panic",
		"go.uber.org/zap.(*Logger).Fatal",
	}
	headerLine := bytes.IndexByte(s, '\n')
	if headerLine == -1 {
		return string(s)
	}
	stack := s[headerLine:]
	targetLine := -1
	for _, target := range targets {
		targetLine = bytes.Index(stack, []byte(target))
		if targetLine != -1 {
			break
		}
	}
	if targetLine == -1 {
		return string(s)
	}
	stack = stack[targetLine+1:]
	// stack has two lines per frame
	for range 2 {
		nextLine := bytes.IndexByte(stack, '\n')
		if nextLine == -1 {
			return string(s)
		}
		stack = stack[nextLine+1:]
	}
	return string(s[:headerLine+1]) + string(stack)
}
