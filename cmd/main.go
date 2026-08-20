package main

import (
	"go.uber.org/zap"
	"net/http"
	"newproject/internal/config"
	"newproject/internal/handlers"
	"newproject/internal/logger"
	"newproject/internal/storage"
)

var cfg *config.Config

func main() {
	cfg = config.New()

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	store, err := storage.New(cfg.DSN)
	if err != nil {
		logger.Log.Fatal("Failed to connect to Database", zap.Error(err))
	}
	defer store.Close()

	h := handlers.New(store, cfg.BaseURL)

	logger.Log.Info("Starting server", zap.String("address", cfg.Addr))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{id}", h.HandleGet)
	mux.HandleFunc("POST /api/shorten", h.HandleAPIShorten)
	mux.HandleFunc("POST /", h.HandleTextShorten)
	if err := http.ListenAndServe(cfg.Addr, logger.RequestLogger(mux)); err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}
