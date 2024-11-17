package http

import (
	"context"
	"fmt"

	"github.com/LLIEPJIOK/weather-forecast/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/pkg/api/weather"

	"github.com/labstack/echo/v4"
)

type WeatherService interface {
	AddWeatherObservation(ctx context.Context, ob models.WeatherObservation) (int, error)
	GetWeatherObservation(ctx context.Context, id int) (models.WeatherObservation, error)
	UpdateWeatherObservation(ctx context.Context, ob models.WeatherObservation) error
	DeleteWeatherObservation(ctx context.Context, id int) (models.WeatherObservation, error)
	ListWeatherObservations(ctx context.Context) ([]models.WeatherObservation, error)
}

type Server struct {
	restServer  *echo.Echo
	restAddress string
}

func New(ctx context.Context, restPort int, weatherService WeatherService) *Server {
	httpSever := echo.New()
	weather.RegisterWeatherObservationRoutes(ctx, httpSever, weatherService)

	return &Server{
		restServer:  httpSever,
		restAddress: fmt.Sprintf(":%d", restPort),
	}
}

func (s *Server) Start() error {
	s.restServer.Logger.Infof("starting rest server at address=%q", s.restAddress)

	if err := s.restServer.Start(s.restAddress); err != nil {
		return fmt.Errorf("failed to start rest server: %w", err)
	}

	return nil
}
