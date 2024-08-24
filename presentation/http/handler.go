package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/wire"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/internal/accesslog"
	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/log"
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
	cs *config.Sentry,
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

	if cs.Enable {
		m, err := newSentryMiddleware(cs)
		if err != nil {
			slog.Error("failed to initialize Sentry middleware", log.Err(err))
		}
		h.middlewares = append(h.middlewares, m.Handle)
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
		slog.InfoContext(ctx, "http: line callback received")

		events, err := h.bot.EventsFromRequest(r)
		if err != nil {
			slog.ErrorContext(ctx, "http: failed to parse request", log.Err(err))
			report(r, "http: failed to parse request", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := h.eventHandler.Handle(ctx, events); err != nil {
			slog.ErrorContext(ctx, "http: failed to handle events", log.Err(err))
			report(r, "http: failed to handle events", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *handler) executeScheduler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.InfoContext(ctx, "http: execute scheduler")

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
			slog.WarnContext(ctx, "http: failed to authorize", log.Err(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if err := h.eventHandler.HandleSchedule(ctx); err != nil {
			slog.ErrorContext(ctx, "http: failed to execute scheduler", log.Err(err))
			report(r, "http: failed to execute scheduler", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handler) executeReminder() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.InfoContext(ctx, "http: execute reminder")

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
			slog.WarnContext(ctx, "http: failed to authorize", log.Err(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(r.Header.Get("content-type"), "application/json") {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		itemIDJSON := new(model.ReminderItemIDJSON)
		if err := json.NewDecoder(r.Body).Decode(itemIDJSON); err != nil {
			slog.WarnContext(ctx, "http: failed to parse request", log.Err(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if itemIDJSON.ConversationID == "" || itemIDJSON.ItemID == "" {
			slog.WarnContext(ctx, "http: invalid payload: conversation_id or item_id is empty",
				slog.String("ConversationID", itemIDJSON.ConversationID),
				slog.String("ItemID", itemIDJSON.ItemID),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := h.eventHandler.HandleReminder(ctx, itemIDJSON); err != nil {
			slog.ErrorContext(ctx, "http: failed to execute reminder", log.Err(err))
			report(r, "http: failed to execute reminder", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handler) serveImage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.InfoContext(ctx, "http: serve image")

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
		sl := slog.With(slog.String("key", key))
		rc, size, err := h.imageHandler.Handle(ctx, key)
		if err != nil {
			if code.From(err) == code.NotFound {
				sl.WarnContext(ctx, "http: image not found", log.Err(err))
				w.WriteHeader(http.StatusNotFound)
				return
			}

			sl.ErrorContext(ctx, "http: failed to serve image", log.Err(err))
			report(r, "http: failed to serve image", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rc.Close()

		w.Header().Set("content-length", strconv.Itoa(size))
		if _, err := io.Copy(w, rc); err != nil {
			if isCanceledByClient(r, err) {
				sl.InfoContext(ctx, "http: request canceled by client",
					slog.Any("errors", []error{err, r.Context().Err()}),
				)
				return
			}

			sl.ErrorContext(ctx, "http: failed to copy image", log.Err(err))
			report(r, "http: failed to copy image", err)
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
