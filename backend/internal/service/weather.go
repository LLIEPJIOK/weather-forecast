package service

import (
	"context"
	"fmt"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
)

//go:generate mockery --name WeatherRepo --structname MockWeatherRepo --filename mock_weather_repo_test.go --outpkg service_test --output .
type WeatherRepo interface {
	AddWeather(ctx context.Context, ob *models.Weather) (int, error)
	GetWeather(ctx context.Context, id int) (*models.Weather, error)
	UpdateWeather(ctx context.Context, ob *models.Weather) error
	DeleteWeather(ctx context.Context, id int) (*models.Weather, error)
	ListWeathers(ctx context.Context) ([]*models.Weather, error)
}

type WeatherService struct {
	repo WeatherRepo
}

func NewWeatherService(repo WeatherRepo) *WeatherService {
	return &WeatherService{repo: repo}
}

func (s *WeatherService) AddWeather(
	ctx context.Context,
	ob *models.Weather,
) (int, error) {
	id, err := s.repo.AddWeather(ctx, ob)
	if err != nil {
		return 0, fmt.Errorf("failed to add weather: %w", err)
	}

	return id, nil
}

func (s *WeatherService) GetWeather(
	ctx context.Context,
	id int,
) (*models.Weather, error) {
	ob, err := s.repo.GetWeather(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: %w", err)
	}

	return ob, nil
}

func (s *WeatherService) UpdateWeather(
	ctx context.Context,
	ob *models.Weather,
) error {
	err := s.repo.UpdateWeather(ctx, ob)
	if err != nil {
		return fmt.Errorf("failed to update weather: %w", err)
	}

	return nil
}

func (s *WeatherService) DeleteWeather(
	ctx context.Context,
	id int,
) (*models.Weather, error) {
	ob, err := s.repo.DeleteWeather(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to delete weather: %w",
			err,
		)
	}

	return ob, nil
}

func (s *WeatherService) ListWeathers(
	ctx context.Context,
) ([]*models.Weather, error) {
	obList, err := s.repo.ListWeathers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list weathers: %w", err)
	}

	return obList, nil
}
