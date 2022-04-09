package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/usecase"
)

// ScreenshotHandlerSet provides a wire set.
var ScreenshotHandlerSet = wire.NewSet(
	NewScreenshotHandler,
	wire.Bind(new(http.Handler), new(*ScreenshotHandler)),
)

type ScreenshotHandler struct {
	log               *logger.Logger
	middlewares       []func(http.Handler) http.Handler
	screenshotHandler usecase.ScreenshotHandler
}

func NewScreenshotHandler(
	log *logger.Logger,
	screenshotHandler usecase.ScreenshotHandler,
) *ScreenshotHandler {
	return &ScreenshotHandler{
		log:               log,
		screenshotHandler: screenshotHandler,
		middlewares: []func(http.Handler) http.Handler{
			XCTCOpenTelemetry(),
			PanicHandler(log),
		},
	}
}

func (h *ScreenshotHandler) registerMiddleware(handler http.Handler) http.Handler {
	for _, middleware := range h.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (h *ScreenshotHandler) healthCheck() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *ScreenshotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.healthCheck())
	mux.HandleFunc("/screenshot", h.screenshot())
	h.registerMiddleware(mux).ServeHTTP(w, r)
}

func (h *ScreenshotHandler) screenshot() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := h.log.WithTraceFromContext(ctx)

		w.Header().Set("content-type", "image/png")
		w.Header().Set("allow", "OPTIONS, HEAD, GET")
		switch r.Method {
		case http.MethodGet:
			// do nothing

		case http.MethodHead, http.MethodOptions:
			w.WriteHeader(http.StatusOK)
			return

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		targetURL := query.Get("url")
		targetSelector := query.Get("selector")

		if targetURL == "" || targetSelector == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Fprintln(w, "url or selector parameter is missing")
			return
		}

		target, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Fprintln(w, "url parameter is invalid")
			return
		}

		img, size, err := h.screenshotHandler.Handle(ctx, target, targetSelector)
		if err != nil {
			cl.Error("failed to screenshot", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-length", strconv.Itoa(size))
		if _, err := io.Copy(w, img); err != nil {
			if isCanceledByClient(r, err) {
				cl.Info("request canceled by client",
					zap.Errors("errors", []error{err, r.Context().Err()}),
				)
				return
			}

			cl.Error("failed to write image", zap.Error(err))
			return
		}
	}
}
