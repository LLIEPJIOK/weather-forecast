package postgres

import (
	"context"
	"fmt"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	PostgresUserName string `env:"POSTGRES_USER"     env-default:"root"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-default:"123"`
	PostgresDBName   string `env:"POSTGRES_DB"       env-default:"weather"`
	PostgresHost     string `env:"POSTGRES_HOST"     env-default:"localhost"`
	PostgresPort     string `env:"PGPORT"            env-default:"5432"`
}

type DB struct {
	db      *sqlx.DB
	queries *Queries
}

func NewPostgres(config PostgresConfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.PostgresUserName,
		config.PostgresPassword,
		config.PostgresDBName,
		config.PostgresHost,
		config.PostgresPort,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if _, err := db.Conn(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{
		queries: New(db),
	}, nil
}

func (db *DB) AddWeather(
	ctx context.Context,
	weather *models.Weather,
) (*models.Weather, error) {
	arg := AddWeatherParams{
		Timestamp:     weather.Timestamp,
		Temperature:   weather.Temperature,
		Humidity:      weather.Humidity,
		Pressure:      weather.Pressure,
		WindSpeed:     weather.WindSpeed,
		City:          weather.City,
		Country:       weather.Country,
		WeatherStatus: weather.WeatherStatus,
	}

	res, err := db.queries.AddWeather(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to add weather: %w", err)
	}

	wth := dbWeatherToGlobal(res)
	return &wth, nil
}

func (db *DB) GetWeather(ctx context.Context, id int) (*models.Weather, error) {
	res, err := db.queries.GetWeather(ctx, int64(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: :%w", err)
	}

	wth := dbWeatherToGlobal(res)
	return &wth, nil
}

func (db *DB) ListWeathers(ctx context.Context) ([]*models.Weather, error) {
	res, err := db.queries.ListWeathers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list weathers: %w", err)
	}

	weathers := make([]*models.Weather, len(res))
	for i, v := range res {
		wth := dbWeatherToGlobal(v)
		weathers[i] = &wth
	}

	return weathers, nil
}

func (db *DB) UpdateWeather(
	ctx context.Context,
	weather *models.Weather,
) (*models.Weather, error) {
	arg := UpdateWeatherParams{
		ID:        int64(weather.ID),
		Timestamp: weather.Timestamp,
		Column3:   weather.Temperature,
		Column4:   weather.Humidity,
		Column5:   weather.Pressure,
		Column6:   weather.WindSpeed,
		Column7:   weather.City,
		Column8:   weather.Country,
		Column9:   weather.WeatherStatus,
	}

	res, err := db.queries.UpdateWeather(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update weather: %w", err)
	}

	ord := dbWeatherToGlobal(res)
	return &ord, nil
}

func (db *DB) DeleteWeather(ctx context.Context, id int) (*models.Weather, error) {
	res, err := db.queries.DeleteWeather(ctx, int64(id))
	if err != nil {
		return nil, fmt.Errorf("failed to delete weather: %w", err)
	}

	ord := dbWeatherToGlobal(res)
	return &ord, nil
}

func (db *DB) Close() error {
	if err := db.db.Close(); err != nil {
		return fmt.Errorf("failed to close postgres: %w", err)
	}

	return nil
}

func dbWeatherToGlobal(weather Weather) models.Weather {
	return models.Weather{
		ID:            int(weather.ID),
		Timestamp:     weather.Timestamp,
		City:          weather.City,
		Country:       weather.Country,
		Temperature:   weather.Temperature,
		Humidity:      weather.Humidity,
		Pressure:      weather.Pressure,
		WindSpeed:     weather.WindSpeed,
		WeatherStatus: weather.WeatherStatus,
	}
}
