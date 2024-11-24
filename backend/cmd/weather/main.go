package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/config"
	"github.com/LLIEPJIOK/weather-forecast/backend/internal/repository"
	"github.com/LLIEPJIOK/weather-forecast/backend/internal/service"
	"github.com/LLIEPJIOK/weather-forecast/backend/internal/tranport/http"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to get config: %s", err)
	}

	whetherRepo := repository.NewWeatherRepository()
	whetherService := service.NewWeatherService(whetherRepo)
	server := http.New(ctx, cfg.RESTServerPort, whetherService)

	if err := server.Start(); err != nil {
		slog.Error(fmt.Sprintf("server.Start(): %s", err))
		os.Exit(1)
	}
}
