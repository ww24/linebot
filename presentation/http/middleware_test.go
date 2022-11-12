package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ww24/linebot/logger"
)

func TestPanicHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		h          http.Handler
		wantStatus int
		wantBody   string
	}{
		{
			name: "success handling",
			h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "success")
			}),
			wantStatus: http.StatusOK,
			wantBody:   "success",
		},
		{
			name: "panic in handler",
			h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("unexpected")
			}),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Internal Server Error\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			panicHandler(logger.NewNop())(tt.h).ServeHTTP(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}
