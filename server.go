package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/loadept/notifier/internal/middleware"
	"github.com/loadept/notifier/internal/webhook"
)

var env = os.Getenv

func main() {
	log.SetFlags(0)
	mux := http.NewServeMux()
	server := &http.Server{
		Handler:      middleware.Logger(mux),
		Addr:         env("ADDR"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	client := &http.Client{Timeout: 10 * time.Second}

	auth := middleware.Auth(env("SECRET_KEY"))
	wh := webhook.NewHandler(client, env("DISCORD_WEBHOOK"))

	mux.Handle("POST /notify", auth(wh.Handler()))

	go func() {
		log.Printf("server listen on %s", server.Addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down server...")
	fatalIfErr(server.Shutdown(shutCtx))
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
