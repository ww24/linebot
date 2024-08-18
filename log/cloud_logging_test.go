package log

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
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
