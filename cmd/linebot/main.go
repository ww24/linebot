package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ww24/linebot/bot"
	"github.com/ww24/linebot/infra/firestore"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	botCfg := bot.Config{
		ChannelSecret:   os.Getenv("LINE_CHANNEL_SECRET"),
		ChannelToken:    os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		ConversationIDs: strings.Split(os.Getenv("ALLOW_CONV_IDS"), ","),
	}
	storeCfg := firestore.ClientConfig{
		ProjectID: os.Getenv("GCP_PROJECT_ID"),
	}
	bot, err := register(ctx, botCfg, storeCfg)
	if err != nil {
		log.Printf("Error: %+v\n", err)
	}

	mux := http.NewServeMux()
	// health check endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		io.WriteString(w, "healthy")
	})
	// line webhook endpoint
	mux.HandleFunc("/line_callback", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received")

		if err := bot.HandleRequest(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Request Error: %+v\n", err)
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
	log.Println("start server")
	//nolint:errcheck
	go srv.ListenAndServe()

	// wait signal
	<-ctx.Done()
	stop()

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		log.Printf("Error: %+v\n", err)
	}
}
