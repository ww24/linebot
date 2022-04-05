package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/usecase"
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewHandler,
	wire.Bind(new(http.Handler), new(*Handler)),
)

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
			XCTCOpenTelemetry(),
			PanicHandler(log),
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.healthCheck())
	mux.HandleFunc("/line_callback", h.lineCallback())
	mux.HandleFunc("/scheduler", h.executeScheduler())
	mux.HandleFunc("/reminder", h.executeReminder())
	h.registerMiddleware(mux).ServeHTTP(w, r)
}

func (h *Handler) registerMiddleware(handler http.Handler) http.Handler {
	for _, middleware := range h.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (h *Handler) healthCheck() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) lineCallback() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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

func (h *Handler) executeScheduler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := h.log.WithTraceFromContext(ctx)
		cl.Info("execute scheduler")

		w.Header().Set("allow", "OPTIONS, HEAD, POST")
		switch r.Method {
		case http.MethodPost:
			// do nothing

		case http.MethodHead, http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
			return

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if err := h.eventHandler.HandleSchedule(ctx); err != nil {
			cl.Error("failed to execute scheduler", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) executeReminder() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := h.log.WithTraceFromContext(ctx)
		cl.Info("execute reminder")

		w.Header().Set("allow", "OPTIONS, HEAD, POST")
		switch r.Method {
		case http.MethodPost:
			// do nothing

		case http.MethodHead, http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
			return

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if !strings.HasPrefix(r.Header.Get("content-type"), "application/json") {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		itemIDJSON := new(model.ReminderItemIDJSON)
		if err := json.NewDecoder(r.Body).Decode(itemIDJSON); err != nil {
			cl.Warn("failed to parse request", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if itemIDJSON.ConversationID == "" || itemIDJSON.ItemID == "" {
			cl.Warn("invalid payload: conversation_id or item_id is empty",
				zap.String("conversation_id", itemIDJSON.ConversationID),
				zap.String("item_id", itemIDJSON.ItemID),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := h.eventHandler.HandleReminder(ctx, itemIDJSON); err != nil {
			cl.Error("failed to execute reminder", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
