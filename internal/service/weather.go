package service

import (
	"context"
	"fmt"

	"github.com/LLIEPJIOK/weather-forecast/internal/models"
)

//go:generate mockery --name WeatherRepo --structname MockWeatherRepo --filename mock_weather_repo_test.go --outpkg service_test --output .
type WeatherRepo interface {
	AddWeatherObservation(ctx context.Context, ob models.WeatherObservation) (int, error)
	GetWeatherObservation(ctx context.Context, id int) (models.WeatherObservation, error)
	UpdateWeatherObservation(ctx context.Context, ob models.WeatherObservation) error
	DeleteWeatherObservation(ctx context.Context, id int) (models.WeatherObservation, error)
	ListWeatherObservations(ctx context.Context) ([]models.WeatherObservation, error)
}

type WeatherService struct {
	repo WeatherRepo
}

func NewWeatherService(repo WeatherRepo) *WeatherService {
	return &WeatherService{repo: repo}
}

func (s *WeatherService) AddWeatherObservation(
	ctx context.Context,
	ob models.WeatherObservation,
) (int, error) {
	id, err := s.repo.AddWeatherObservation(ctx, ob)
	if err != nil {
		return 0, fmt.Errorf("failed to add weather observation: %w", err)
	}

	return id, nil
}

func (s *WeatherService) GetWeatherObservation(
	ctx context.Context,
	id int,
) (models.WeatherObservation, error) {
	ob, err := s.repo.GetWeatherObservation(ctx, id)
	if err != nil {
		return models.WeatherObservation{}, fmt.Errorf("failed to get weather observation: %w", err)
	}

	return ob, nil
}

func (s *WeatherService) UpdateWeatherObservation(
	ctx context.Context,
	ob models.WeatherObservation,
) error {
	err := s.repo.UpdateWeatherObservation(ctx, ob)
	if err != nil {
		return fmt.Errorf("failed to update weather observation: %w", err)
	}

	return nil
}

func (s *WeatherService) DeleteWeatherObservation(
	ctx context.Context,
	id int,
) (models.WeatherObservation, error) {
	ob, err := s.repo.DeleteWeatherObservation(ctx, id)
	if err != nil {
		return models.WeatherObservation{}, fmt.Errorf(
			"failed to delete weather observation: %w",
			err,
		)
	}

	return ob, nil
}

func (s *WeatherService) ListWeatherObservations(
	ctx context.Context,
) ([]models.WeatherObservation, error) {
	obList, err := s.repo.ListWeatherObservations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list weather observations: %w", err)
	}

	return obList, nil
}
