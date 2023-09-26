package logger

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ww24/linebot/internal/buildinfo"
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

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceReference
type sourceReference struct {
	Repository string
	RevisionID string
}

func (r *sourceReference) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("repository", strings.Replace(r.Repository, "git://", "https://", 1))
	enc.AddString("revisionId", r.RevisionID)
	return nil
}

type sourceReferences []*sourceReference

func (r sourceReferences) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, sr := range r {
		if err := enc.AppendObject(sr); err != nil {
			return fmt.Errorf("failed to append sourceReference: %w", err)
		}
	}
	return nil
}

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext
type errorContext struct {
	ReportLocation   *sourceLocation
	SourceReferences sourceReferences
}

func (c *errorContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if c.ReportLocation != nil {
		if err := enc.AddObject("reportLocation", c.ReportLocation); err != nil {
			return fmt.Errorf("failed to add reportLocation field: %w", err)
		}
	}
	if len(c.SourceReferences) > 0 {
		if err := enc.AddArray("sourceReferences", c.SourceReferences); err != nil {
			return fmt.Errorf("failed to add sourceReferences field: %w", err)
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
		SourceReferences: sourceReferences{
			{
				Repository: buildinfo.Repository(),
				RevisionID: buildinfo.Revision(),
			},
		},
	}
}

// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
type logEntrySourceLocation struct {
	File     string
	Line     int
	Function string
}

func (l *logEntrySourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("file", l.File)
	enc.AddString("line", strconv.Itoa(l.Line))
	enc.AddString("function", l.Function)
	return nil
}

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
