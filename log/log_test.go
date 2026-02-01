package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_chopStack(t *testing.T) {
	t.Parallel()
	type args struct {
		s      []byte
		target string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				s:      []byte{},
				target: "github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func5",
			},
			want: "",
		},
		{
			name: "header only",
			args: args{
				s:      []byte("goroutine 1 [running]:\n"),
				target: "github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func5",
			},
			want: "goroutine 1 [running]:\n",
		},
		{
			name: "target frame not found",
			args: args{
				s: []byte(`goroutine 120 [running]:
runtime/debug.Stack()
	runtime/debug/stack.go:24 +0x5e
github.com/ww24/linebot/logger.(*core).Write(0xc008d3fee0, {0x2, {0xc13ce36b6e3319a3, 0x199bb75c0e, 0x2a21fcd3dc60}, {0x0, 0x0}, {0x2a21fb10fda0, 0x1d}, {0x1, ...}, ...}, ...)
	github.com/ww24/linebot/logger/core.go:41 +0x2f8
go.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc0001108f0, {0xc00940ee80, 0x1, 0x1})
	go.uber.org/zap@v1.26.0/zapcore/entry.go:253 +0x1dc
go.uber.org/zap.(*Logger).Warn(0x2a21fc46cf40?, {0x2a21fb10fda0?, 0x2a21fab10566?}, {0xc00940ee80, 0x1, 0x1})
	go.uber.org/zap@v1.26.0/logger.go:254 +0x51
github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func5({0x2a21fc4774f0, 0xc0093ff6c0}, 0x0?)
	github.com/ww24/linebot/presentation/http/handler.go:89 +0x1b0
net/http.HandlerFunc.ServeHTTP(0x20?, {0x2a21fc4774f0?, 0xc0093ff6c0?}, 0xc00809b278?)
	net/http/server.go:2136 +0x29
net/http.(*ServeMux).ServeHTTP(0x2a21fc478710?, {0x2a21fc4774f0, 0xc0093ff6c0}, 0xc008cdfe00)
	net/http/server.go:2514 +0x142
github.com/ww24/linebot/presentation/http.NewHandler.accessLogHandler.func3.1({0x2a21fc478710, 0xc0094051a0}, 0xc008cdfe00)
	github.com/ww24/linebot/presentation/http/middleware.go:72 +0x75f
net/http.HandlerFunc.ServeHTTP(0x2a21fc479120?, {0x2a21fc478710?, 0xc0094051a0?}, 0x2a21fbd78f90?)
	net/http/server.go:2136 +0x29
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.(*middleware).serveHTTP(0xc0004de4d0, {0x2a21fc4771f0?, 0xc00927ed20}, 0xc008cdfc00, {0x2a21fc46add8, 0xc0000331d0})
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:217 +0x1202
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.NewMiddleware.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0x10?)
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:81 +0x35
net/http.HandlerFunc.ServeHTTP(0x3eb1554eeb98?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x2a21fa8abe74?)
	net/http/server.go:2136 +0x29
github.com/ww24/linebot/presentation/http.NewHandler.panicHandler.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0xc009409ce0?)
	github.com/ww24/linebot/presentation/http/middleware.go:33 +0x78
net/http.HandlerFunc.ServeHTTP(0x2a21fcd6d380?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0xc000093b50?)
	net/http/server.go:2136 +0x29
net/http.serverHandler.ServeHTTP({0xc009458780?}, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x6?)
	net/http/server.go:2938 +0x8e
net/http.(*conn).serve(0xc0085bc6c0, {0x2a21fc479120, 0xc0014ad260})
	net/http/server.go:2009 +0x5f4`),
				target: "github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func4",
			},
			want: `goroutine 120 [running]:
runtime/debug.Stack()
	runtime/debug/stack.go:24 +0x5e
github.com/ww24/linebot/logger.(*core).Write(0xc008d3fee0, {0x2, {0xc13ce36b6e3319a3, 0x199bb75c0e, 0x2a21fcd3dc60}, {0x0, 0x0}, {0x2a21fb10fda0, 0x1d}, {0x1, ...}, ...}, ...)
	github.com/ww24/linebot/logger/core.go:41 +0x2f8
go.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc0001108f0, {0xc00940ee80, 0x1, 0x1})
	go.uber.org/zap@v1.26.0/zapcore/entry.go:253 +0x1dc
go.uber.org/zap.(*Logger).Warn(0x2a21fc46cf40?, {0x2a21fb10fda0?, 0x2a21fab10566?}, {0xc00940ee80, 0x1, 0x1})
	go.uber.org/zap@v1.26.0/logger.go:254 +0x51
github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func5({0x2a21fc4774f0, 0xc0093ff6c0}, 0x0?)
	github.com/ww24/linebot/presentation/http/handler.go:89 +0x1b0
net/http.HandlerFunc.ServeHTTP(0x20?, {0x2a21fc4774f0?, 0xc0093ff6c0?}, 0xc00809b278?)
	net/http/server.go:2136 +0x29
net/http.(*ServeMux).ServeHTTP(0x2a21fc478710?, {0x2a21fc4774f0, 0xc0093ff6c0}, 0xc008cdfe00)
	net/http/server.go:2514 +0x142
github.com/ww24/linebot/presentation/http.NewHandler.accessLogHandler.func3.1({0x2a21fc478710, 0xc0094051a0}, 0xc008cdfe00)
	github.com/ww24/linebot/presentation/http/middleware.go:72 +0x75f
net/http.HandlerFunc.ServeHTTP(0x2a21fc479120?, {0x2a21fc478710?, 0xc0094051a0?}, 0x2a21fbd78f90?)
	net/http/server.go:2136 +0x29
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.(*middleware).serveHTTP(0xc0004de4d0, {0x2a21fc4771f0?, 0xc00927ed20}, 0xc008cdfc00, {0x2a21fc46add8, 0xc0000331d0})
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:217 +0x1202
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.NewMiddleware.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0x10?)
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:81 +0x35
net/http.HandlerFunc.ServeHTTP(0x3eb1554eeb98?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x2a21fa8abe74?)
	net/http/server.go:2136 +0x29
github.com/ww24/linebot/presentation/http.NewHandler.panicHandler.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0xc009409ce0?)
	github.com/ww24/linebot/presentation/http/middleware.go:33 +0x78
net/http.HandlerFunc.ServeHTTP(0x2a21fcd6d380?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0xc000093b50?)
	net/http/server.go:2136 +0x29
net/http.serverHandler.ServeHTTP({0xc009458780?}, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x6?)
	net/http/server.go:2938 +0x8e
net/http.(*conn).serve(0xc0085bc6c0, {0x2a21fc479120, 0xc0014ad260})
	net/http/server.go:2009 +0x5f4`,
		},
		{
			name: "target frame found",
			args: args{
				s: []byte(`goroutine 120 [running]:
runtime/debug.Stack()
	runtime/debug/stack.go:24 +0x5e
github.com/ww24/linebot/logger.(*core).Write(0xc008d3fee0, {0x2, {0xc13ce36b6e3319a3, 0x199bb75c0e, 0x2a21fcd3dc60}, {0x0, 0x0}, {0x2a21fb10fda0, 0x1d}, {0x1, ...}, ...}, ...)
	github.com/ww24/linebot/logger/core.go:41 +0x2f8
go.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc0001108f0, {0xc00940ee80, 0x1, 0x1})
	go.uber.org/zap@v1.26.0/zapcore/entry.go:253 +0x1dc
go.uber.org/zap.(*Logger).Error(0x2a21fc46cf40?, {0x2a21fb10fda0?, 0x2a21fab10566?}, {0xc00940ee80, 0x1, 0x1})
	go.uber.org/zap@v1.26.0/logger.go:262 +0x51
github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func5({0x2a21fc4774f0, 0xc0093ff6c0}, 0x0?)
	github.com/ww24/linebot/presentation/http/handler.go:89 +0x1b0
net/http.HandlerFunc.ServeHTTP(0x20?, {0x2a21fc4774f0?, 0xc0093ff6c0?}, 0xc00809b278?)
	net/http/server.go:2136 +0x29
net/http.(*ServeMux).ServeHTTP(0x2a21fc478710?, {0x2a21fc4774f0, 0xc0093ff6c0}, 0xc008cdfe00)
	net/http/server.go:2514 +0x142
github.com/ww24/linebot/presentation/http.NewHandler.accessLogHandler.func3.1({0x2a21fc478710, 0xc0094051a0}, 0xc008cdfe00)
	github.com/ww24/linebot/presentation/http/middleware.go:72 +0x75f
net/http.HandlerFunc.ServeHTTP(0x2a21fc479120?, {0x2a21fc478710?, 0xc0094051a0?}, 0x2a21fbd78f90?)
	net/http/server.go:2136 +0x29
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.(*middleware).serveHTTP(0xc0004de4d0, {0x2a21fc4771f0?, 0xc00927ed20}, 0xc008cdfc00, {0x2a21fc46add8, 0xc0000331d0})
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:217 +0x1202
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.NewMiddleware.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0x10?)
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:81 +0x35
net/http.HandlerFunc.ServeHTTP(0x3eb1554eeb98?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x2a21fa8abe74?)
	net/http/server.go:2136 +0x29
github.com/ww24/linebot/presentation/http.NewHandler.panicHandler.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0xc009409ce0?)
	github.com/ww24/linebot/presentation/http/middleware.go:33 +0x78
net/http.HandlerFunc.ServeHTTP(0x2a21fcd6d380?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0xc000093b50?)
	net/http/server.go:2136 +0x29
net/http.serverHandler.ServeHTTP({0xc009458780?}, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x6?)
	net/http/server.go:2938 +0x8e
net/http.(*conn).serve(0xc0085bc6c0, {0x2a21fc479120, 0xc0014ad260})
	net/http/server.go:2009 +0x5f4`),
				target: "github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func5",
			},
			want: `goroutine 120 [running]:
github.com/ww24/linebot/presentation/http.NewHandler.(*handler).lineCallback.func5({0x2a21fc4774f0, 0xc0093ff6c0}, 0x0?)
	github.com/ww24/linebot/presentation/http/handler.go:89 +0x1b0
net/http.HandlerFunc.ServeHTTP(0x20?, {0x2a21fc4774f0?, 0xc0093ff6c0?}, 0xc00809b278?)
	net/http/server.go:2136 +0x29
net/http.(*ServeMux).ServeHTTP(0x2a21fc478710?, {0x2a21fc4774f0, 0xc0093ff6c0}, 0xc008cdfe00)
	net/http/server.go:2514 +0x142
github.com/ww24/linebot/presentation/http.NewHandler.accessLogHandler.func3.1({0x2a21fc478710, 0xc0094051a0}, 0xc008cdfe00)
	github.com/ww24/linebot/presentation/http/middleware.go:72 +0x75f
net/http.HandlerFunc.ServeHTTP(0x2a21fc479120?, {0x2a21fc478710?, 0xc0094051a0?}, 0x2a21fbd78f90?)
	net/http/server.go:2136 +0x29
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.(*middleware).serveHTTP(0xc0004de4d0, {0x2a21fc4771f0?, 0xc00927ed20}, 0xc008cdfc00, {0x2a21fc46add8, 0xc0000331d0})
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:217 +0x1202
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp.NewMiddleware.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0x10?)
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.44.0/handler.go:81 +0x35
net/http.HandlerFunc.ServeHTTP(0x3eb1554eeb98?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x2a21fa8abe74?)
	net/http/server.go:2136 +0x29
github.com/ww24/linebot/presentation/http.NewHandler.panicHandler.func1.1({0x2a21fc4771f0?, 0xc00927ed20?}, 0xc009409ce0?)
	github.com/ww24/linebot/presentation/http/middleware.go:33 +0x78
net/http.HandlerFunc.ServeHTTP(0x2a21fcd6d380?, {0x2a21fc4771f0?, 0xc00927ed20?}, 0xc000093b50?)
	net/http/server.go:2136 +0x29
net/http.serverHandler.ServeHTTP({0xc009458780?}, {0x2a21fc4771f0?, 0xc00927ed20?}, 0x6?)
	net/http/server.go:2938 +0x8e
net/http.(*conn).serve(0xc0085bc6c0, {0x2a21fc479120, 0xc0014ad260})
	net/http/server.go:2009 +0x5f4`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := chopStack(tt.args.s, tt.args.target)
			assert.Equal(t, tt.want, got)
		})
	}
}
