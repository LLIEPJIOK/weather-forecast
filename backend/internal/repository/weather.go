package repository

import (
	"context"
	"fmt"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
)

//go:generate mockery --name Database --structname MockDatabase --filename mock_database_test.go --outpkg repository_test --output .
type Database interface {
	AddWeather(ctx context.Context, weather *models.Weather) (*models.Weather, error)
	GetWeather(ctx context.Context, id int) (*models.Weather, error)
	ListWeathers(ctx context.Context) ([]*models.Weather, error)
	UpdateWeather(ctx context.Context, weather *models.Weather) (*models.Weather, error)
	DeleteWeather(ctx context.Context, id int) (*models.Weather, error)
}

type WeatherRepository struct {
	db Database
}

func NewWeatherRepository(db Database) *WeatherRepository {
	return &WeatherRepository{
		db: db,
	}
}

func (r *WeatherRepository) AddWeather(
	ctx context.Context,
	ob *models.Weather,
) (int, error) {
	res, err := r.db.AddWeather(ctx, ob)
	if err != nil {
		return 0, fmt.Errorf("failed to add weather: %w", err)
	}

	return res.ID, nil
}

func (r *WeatherRepository) GetWeather(
	ctx context.Context,
	id int,
) (*models.Weather, error) {
	res, err := r.db.GetWeather(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: %w", err)
	}

	return res, nil
}

func (r *WeatherRepository) UpdateWeather(
	ctx context.Context,
	ob *models.Weather,
) error {
	_, err := r.db.UpdateWeather(ctx, ob)
	if err != nil {
		return fmt.Errorf("failed to update weather: %w", err)
	}

	return nil
}

func (r *WeatherRepository) DeleteWeather(
	ctx context.Context,
	id int,
) (*models.Weather, error) {
	res, err := r.db.DeleteWeather(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete weather: %w", err)
	}

	return res, nil
}

func (r *WeatherRepository) ListWeathers(
	ctx context.Context,
) ([]*models.Weather, error) {
	res, err := r.db.ListWeathers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list weathers: %w", err)
	}

	return res, nil
}
