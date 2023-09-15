package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/internal/accesslog"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/tracer"
	"github.com/ww24/linebot/usecase"
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewHandler,
	NewAuthorizer,
)

type handler struct {
	bot          service.Bot
	auth         *Authorizer
	eventHandler usecase.EventHandler
	imageHandler usecase.ImageHandler
	middlewares  []func(http.Handler) http.Handler
}

func NewHandler(
	bot service.Bot,
	auth *Authorizer,
	eventHandler usecase.EventHandler,
	imageHandler usecase.ImageHandler,
	publisher accesslog.Publisher,
	cfg *config.AccessLog,
) (http.Handler, error) {
	h := &handler{
		bot:          bot,
		auth:         auth,
		eventHandler: eventHandler,
		imageHandler: imageHandler,
		middlewares: []func(http.Handler) http.Handler{
			panicHandler(),
			tracer.HTTPMiddleware(),
			accessLogHandler(publisher, cfg),
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.healthCheck())
	mux.HandleFunc("/line_callback", h.lineCallback())
	mux.HandleFunc("/scheduler", h.executeScheduler())
	mux.HandleFunc("/reminder", h.executeReminder())
	mux.HandleFunc("/image/", h.serveImage())
	return h.registerMiddleware(mux), nil
}

func (h *handler) registerMiddleware(handler http.Handler) http.Handler {
	for i := len(h.middlewares) - 1; i >= 0; i-- {
		handler = h.middlewares[i](handler)
	}
	return handler
}

func (h *handler) healthCheck() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *handler) lineCallback() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := logger.Default(ctx)
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

func (h *handler) executeScheduler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := logger.Default(ctx)
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

		if err := h.auth.Authorize(ctx, r); err != nil {
			cl.Info("failed to authorize", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

func (h *handler) executeReminder() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := logger.Default(ctx)
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

		if err := h.auth.Authorize(ctx, r); err != nil {
			cl.Info("failed to authorize", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

func (h *handler) serveImage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := logger.Default(ctx)
		cl.Info("serve image")

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

		const prefix = "/image/"
		if !strings.HasPrefix(r.URL.Path, prefix+"weather/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		key := strings.TrimPrefix(r.URL.Path, prefix)
		rc, size, err := h.imageHandler.Handle(ctx, key)
		if err != nil {
			cl.Error("failed to serve image", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rc.Close()

		w.Header().Set("content-length", strconv.Itoa(size))
		if _, err := io.Copy(w, rc); err != nil {
			if isCanceledByClient(r, err) {
				cl.Info("request canceled by client",
					zap.Errors("errors", []error{err, r.Context().Err()}),
				)
				return
			}

			cl.Error("failed to copy image", zap.Error(err))
			return
		}
	}
}

func isCanceledByClient(r *http.Request, err error) bool {
	if err == nil {
		return false
	}
	return (errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, context.Canceled)) &&
		errors.Is(r.Context().Err(), context.Canceled)
}
