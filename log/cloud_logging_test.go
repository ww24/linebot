package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
	"testing/slogtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_severity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level slog.Level
		want  string
	}{
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelDebug + 1, "DEBUG"},
		{slog.LevelDebug + 2, "DEBUG"},
		{slog.LevelDebug + 3, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelInfo + 1, "NOTICE"},
		{slog.LevelInfo + 2, "NOTICE"},
		{slog.LevelInfo + 3, "NOTICE"},
		{slog.LevelWarn, "WARNING"},
		{slog.LevelWarn + 1, "WARNING"},
		{slog.LevelWarn + 2, "WARNING"},
		{slog.LevelWarn + 3, "WARNING"},
		{slog.LevelError, "ERROR"},
		{slog.LevelError + 1, "CRITICAL"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := severity(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCloudLoggingHandler(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	dec := json.NewDecoder(buf)
	slogtest.Run(t, func(t *testing.T) slog.Handler {
		return newCloudLoggingHandler(buf, "project")
	}, func(t *testing.T) map[string]any {
		var m map[string]any
		require.NoError(t, dec.Decode(&m))
		if v, ok := m["timestamp"]; ok {
			m[slog.TimeKey] = v
		}
		if v, ok := m["severity"]; ok {
			m[slog.LevelKey] = v
		}
		if v, ok := m["message"]; ok {
			m[slog.MessageKey] = v
		}
		fmt.Println(m)
		return m
	})
}
