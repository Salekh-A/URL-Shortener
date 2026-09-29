package main

import (
	"newproject/internal/app"
	"newproject/internal/config"
	"newproject/internal/logger"
)

func main() {
	cfg := config.New()

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	application, err := app.New(cfg)
	if err != nil {
		logger.Log.Fatal("Application initialization failed")
	}

	defer application.Close()

	application.Run()
}
