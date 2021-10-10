package http

import (
	"net/http"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"

	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/usecase"
)

// Set provides a wire set.
//nolint: gochecknoglobals
var Set = wire.NewSet(
	NewHandler,
	wire.Bind(new(http.Handler), new(*Handler)),
)

//nolint: gochecknoglobals
var prop = propagator.New()

type Handler struct {
	log          *logger.Logger
	bot          service.Bot
	eventHandler usecase.EventHandler
	middlewares  []func(http.Handler) http.Handler
}

func NewHandler(
	log *logger.Logger,
	bot service.Bot,
	eventHandler usecase.EventHandler,
) *Handler {
	return &Handler{
		log:          log,
		bot:          bot,
		eventHandler: eventHandler,
		middlewares: []func(http.Handler) http.Handler{
			PanicHandler(log),
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/line_callback", h.lineCallback())
	h.registerMiddleware(mux).ServeHTTP(w, r)
}

func (h *Handler) registerMiddleware(handler http.Handler) http.Handler {
	for _, middleware := range h.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (h *Handler) lineCallback() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		cl := h.log.WithTraceFromContext(ctx)
		cl.Info("line callback received")

		events, err := h.bot.EventsFromRequest(r)
		if err != nil {
			cl.Info("failed to parse request", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := h.eventHandler.Handle(ctx, events); err != nil {
			cl.Error("failed to handle events", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
