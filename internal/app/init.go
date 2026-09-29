package app

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"newproject/internal/adapters/pgrepo"
	redisadapter "newproject/internal/adapters/redis"
	"newproject/internal/config"
	"newproject/internal/handlers"
	"newproject/internal/logger"
	"newproject/internal/services"
	"newproject/pkg"
)

type App struct {
	server *http.Server
	db     *pgxpool.Pool
	redis  *redis.Client
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	db, err := pkg.NewDB(ctx)
	if err != nil {
		return nil, err
	}

	redisClient, err := pkg.NewRedis(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	repository := pgrepo.New(db)
	cache := redisadapter.New(redisClient)

	service := services.NewURLService(repository, cache)
	handler := handlers.New(service, cfg.BaseURL)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", handler.HandleGet)
	mux.HandleFunc("POST /api/shorten", handler.HandleAPIShorten)
	mux.HandleFunc("POST /", handler.HandleTextShorten)
	mux.HandleFunc("POST /api/shorten/batch", handler.BatchShorten)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: logger.RequestLogger(mux),
	}

	return &App{
		server: server,
		db:     db,
		redis:  redisClient,
	}, nil
}

func (a *App) Run() {
	logger.Log.Info(
		"Starting server",
		zap.String("address", a.server.Addr),
	)

	go func() {
		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal(
				"Server failed",
				zap.Error(err),
			)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	logger.Log.Info("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error(
			"Server shutdown failed",
			zap.Error(err),
		)
	}

	logger.Log.Info("Server stopped")
}

func (a *App) Close() {
	if a.redis != nil {
		a.redis.Close()
	}

	if a.db != nil {
		a.db.Close()
	}
}
