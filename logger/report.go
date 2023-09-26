package logger

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ServiceContext
// see: https://cloud.google.com/error-reporting/docs/formatting-error-messages
type serviceContext struct {
	service string
	version string
}

func newServiceContext(service, version string) zap.Field {
	sc := &serviceContext{
		service: service,
		version: version,
	}
	return zap.Object("serviceContext", sc)
}

func (c *serviceContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("service", c.service)
	if c.version != "" {
		enc.AddString("version", c.version)
	}
	return nil
}

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
type sourceLocation struct {
	FilePath     string
	LineNumber   int
	FunctionName string
}

func (l *sourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("filePath", l.FilePath)
	enc.AddInt("lineNumber", l.LineNumber)
	enc.AddString("functionName", l.FunctionName)
	return nil
}

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext
type errorContext struct {
	ReportLocation *sourceLocation
}

func (c *errorContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if c.ReportLocation != nil {
		if err := enc.AddObject("reportLocation", c.ReportLocation); err != nil {
			return fmt.Errorf("failed to add reportLocation field: %w", err)
		}
	}
	return nil
}

func newErrorReport(caller zapcore.EntryCaller) *errorContext {
	if !caller.Defined {
		return nil
	}
	return &errorContext{
		ReportLocation: &sourceLocation{
			FilePath:     caller.File,
			LineNumber:   caller.Line,
			FunctionName: caller.Function,
		},
	}
}

// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
//
//nolint:unused
type logEntrySourceLocation struct {
	File     string
	Line     int
	Function string
}

//nolint:unused
func (l *logEntrySourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("file", l.File)
	enc.AddString("line", strconv.Itoa(l.Line))
	enc.AddString("function", l.Function)
	return nil
}

//nolint:unused
func newSourceLocation(caller zapcore.EntryCaller) *logEntrySourceLocation {
	if !caller.Defined {
		return nil
	}
	return &logEntrySourceLocation{
		File:     caller.File,
		Line:     caller.Line,
		Function: caller.Function,
	}
}
