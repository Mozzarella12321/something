package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mozzarella12321/orders-api/internal/config"
	"github.com/mozzarella12321/orders-api/internal/http-server/handlers/url/save"
	"github.com/mozzarella12321/orders-api/internal/http-server/middleware/logger"
	"github.com/mozzarella12321/orders-api/internal/lib/logger/sl"
	"github.com/mozzarella12321/orders-api/internal/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg.StoragePath) //delete
	log := setupLogger(cfg.Env)
	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	storage, err := postgresql.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
	//TODO: run server:
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

// err = storage.SaveURL("google.com", "google")
// if err != nil {
// 	log.Error("failed to save url", sl.Err(err))
// 	os.Exit(1)
// }
// log.Info("saved url", slog.String("url", "google.com"))

// urlToGet, err := storage.GetURL("google")
// if err != nil {
// 	log.Error("failed to get url", sl.Err(err))
// 	os.Exit(1)
// }
// log.Info("got url", slog.String("url", urlToGet))

// urlToDelete, err := storage.DeleteUrl("google")
// if err != nil {
// 	log.Error("failed to delete url", sl.Err(err))
// 	os.Exit(1)
// }
// log.Info("deleted url", slog.String("url", urlToDelete))
