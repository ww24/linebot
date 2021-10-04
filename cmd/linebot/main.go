package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/profiler"

	"github.com/ww24/linebot/bot"
)

const (
	serviceName = "linebot"
)

var (
	// version is set during build
	version string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	botCfg := bot.Config{
		ChannelSecret:   os.Getenv("LINE_CHANNEL_SECRET"),
		ChannelToken:    os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		ConversationIDs: strings.Split(os.Getenv("ALLOW_CONV_IDS"), ","),
	}
	bot, err := register(ctx, botCfg)
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	// initialize cloud profiler if build is production
	if version != "" {
		profilerConfig := profiler.Config{
			Service:           serviceName,
			ServiceVersion:    version,
			MutexProfiling:    true,
			EnableOCTelemetry: true,
		}
		if err := profiler.Start(profilerConfig); err != nil {
			// just log the error if failed to initialize profiler
			log.Printf("failed to start cloud profiler: %+v", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/line_callback", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received")

		if err := bot.HandleRequest(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Request Error: %+v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	addr := ":8000"
	if a := os.Getenv("PORT"); a != "" {
		addr = ":" + a
	}
	srv := &http.Server{
		Handler: mux,
		Addr:    addr,
	}
	log.Printf("start server: %s", version)
	//nolint:errcheck
	go srv.ListenAndServe()

	// wait signal
	<-ctx.Done()
	stop()

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		log.Printf("Error: %+v", err)
	}
}
