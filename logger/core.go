package logger

import (
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
			fields = append(fields, zap.Object("context", report))
		}
	}
	if loc := newSourceLocation(e.Caller); loc != nil {
		fields = append(fields, zap.Object("logging.googleapis.com/sourceLocation", loc))
	}

	//nolint: wrapcheck
	return c.Core.Write(e, fields)
}
