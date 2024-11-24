package http

import (
	"context"
	"fmt"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/backend/pkg/api/weather"

	"github.com/labstack/echo/v4"
)

type WeatherService interface {
	AddWeather(ctx context.Context, ob *models.Weather) (int, error)
	GetWeather(ctx context.Context, id int) (*models.Weather, error)
	UpdateWeather(ctx context.Context, ob *models.Weather) error
	DeleteWeather(ctx context.Context, id int) (*models.Weather, error)
	ListWeathers(ctx context.Context) ([]*models.Weather, error)
}

type Server struct {
	restServer  *echo.Echo
	restAddress string
}

func New(ctx context.Context, restPort int, weatherService WeatherService) *Server {
	httpSever := echo.New()
	weather.RegisterWeatherRoutes(ctx, httpSever, weatherService)

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
