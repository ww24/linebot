package http

import (
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsCanceledByClient(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		requestBody    io.Reader
		requestTimeout time.Duration
		handler        func(*setter) http.HandlerFunc
		want           bool
	}{
		{
			name:           "long task is canceled by client",
			requestBody:    http.NoBody,
			requestTimeout: time.Second,
			handler: func(result *setter) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					err := doLongTask(ctx, 5*time.Second)
					got := isCanceledByClient(r, err)
					result.Set(got)
				}
			},
			want: true, // error is context.Canceled
		},
		{
			name:           "write large data task is canceled by client",
			requestBody:    http.NoBody,
			requestTimeout: time.Second,
			handler: func(result *setter) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					_, err := io.Copy(w, rand.Reader)
					got := isCanceledByClient(r, err)
					result.Set(got)
				}
			},
			want: true, // error is syscall.EPIPE on macOS, syscall.ECONNRESET on Linux
		},
		{
			name:           "read large data task is canceled by client",
			requestBody:    rand.Reader,
			requestTimeout: time.Second,
			handler: func(result *setter) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					_, err := io.Copy(io.Discard, r.Body)
					got := isCanceledByClient(r, err)
					result.Set(got)
				}
			},
			want: true, // error is io.ErrUnexpectedEOF
		},
		{
			name:           "success",
			requestBody:    http.NoBody,
			requestTimeout: 2 * time.Second,
			handler: func(result *setter) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					err := doLongTask(ctx, 100*time.Millisecond)
					got := isCanceledByClient(r, err)
					result.Set(got)
				}
			},
			want: false, // error is nil
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := newSetter(1)
			ts := httptest.NewServer(tt.handler(result))
			defer ts.Close()

			ctx, cancel := context.WithTimeout(context.Background(), tt.requestTimeout)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, tt.requestBody)
			require.NoError(t, err)
			res, err := http.DefaultClient.Do(req)
			if err == nil {
				defer res.Body.Close()
			}

			assert.Equal(t, tt.want, result.Get())
		})
	}
}

type setter struct {
	mutex *sync.Mutex
	wg    *sync.WaitGroup
	value bool
}

func newSetter(count int) *setter {
	wg := new(sync.WaitGroup)
	wg.Add(count)
	return &setter{
		mutex: new(sync.Mutex),
		wg:    wg,
		value: false,
	}
}

func (s *setter) Set(v bool) {
	defer s.wg.Done()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = v
}

func (s *setter) Get() bool {
	s.wg.Wait()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value
}

func doLongTask(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}

	return ctx.Err()
}
