package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"

	del "rest_url_shorter/internal/http-server/handlers/url/delete"
	"rest_url_shorter/internal/http-server/handlers/url/redirect"
	"rest_url_shorter/internal/http-server/handlers/url/save"
	mwlogger "rest_url_shorter/internal/http-server/middleware/logger"
	"rest_url_shorter/internal/storage/postgres"

	"rest_url_shorter/internal/config"
	"rest_url_shorter/internal/lib/logger/sl"
	_ "rest_url_shorter/migrations"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting server",
		slog.String("address", cfg.HTTPServer.Address),
		slog.String("env", cfg.Env),
	)

	log.Debug("Debug messages are enabled")

	storage, err := postgres.New(cfg.ConnectionDB)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	log.Debug("Storage started")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortner", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", del.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
