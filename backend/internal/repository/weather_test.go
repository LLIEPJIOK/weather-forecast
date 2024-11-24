package repository_test

import (
	"context"
	"testing"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/backend/internal/repository"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type databaseBuilder func(t *testing.T) repository.Database

func TestAddWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		ctx       context.Context
		dbBuilder databaseBuilder
		ob        *models.Weather
	}{
		{
			name: "Add first weather observation with positive values",
			ctx:  context.Background(),
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(&models.Weather{
						ID:            1,
						Temperature:   25.0,
						Humidity:      60.0,
						Pressure:      1013.0,
						WeatherStatus: "Clear",
						City:          "Berlin",
						Country:       "Germany",
					}, nil).
					Once()
				mockService.
					On("GetWeather", mock.Anything, 1).
					Return(&models.Weather{
						ID:            1,
						Temperature:   25.0,
						Humidity:      60.0,
						Pressure:      1013.0,
						WeatherStatus: "Clear",
						City:          "Berlin",
						Country:       "Germany",
					}, nil).
					Once()

				return mockService
			},
			ob: &models.Weather{
				ID:            1,
				Temperature:   25.0,
				Humidity:      60.0,
				Pressure:      1013.0,
				WeatherStatus: "Clear",
				City:          "Berlin",
				Country:       "Germany",
			},
		},
		{
			name: "Add second weather observation with different location",
			ctx:  context.Background(),
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(&models.Weather{
						ID:            1,
						Temperature:   30.0,
						Humidity:      65.0,
						Pressure:      1015.0,
						WeatherStatus: "Sunny",
						City:          "Paris",
						Country:       "France",
					}, nil).
					Once()
				mockService.
					On("GetWeather", mock.Anything, 1).
					Return(&models.Weather{
						ID:            1,
						Temperature:   30.0,
						Humidity:      65.0,
						Pressure:      1015.0,
						WeatherStatus: "Sunny",
						City:          "Paris",
						Country:       "France",
					}, nil).
					Once()

				return mockService
			},
			ob: &models.Weather{
				Temperature:   30.0,
				Humidity:      65.0,
				Pressure:      1015.0,
				WeatherStatus: "Sunny",
				City:          "Paris",
				Country:       "France",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.dbBuilder(t)
			repo := repository.NewWeatherRepository(db)

			id, err := repo.AddWeather(tc.ctx, tc.ob)
			require.NoError(t, err)

			tc.ob.ID = id

			ob, err := repo.GetWeather(tc.ctx, id)
			require.NoError(t, err)

			assert.Equal(t, tc.ob, ob)
		})
	}
}

func TestGetWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		ctx       context.Context
		dbBuilder databaseBuilder
		id        int
	}{
		{
			name: "Attempt to get non-existent weather observation",
			ctx:  context.Background(),
			id:   1234,
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("GetWeather", mock.Anything, 1234).
					Return(nil, repository.NewErrNotFound(1234)).
					Once()

				return mockService
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.dbBuilder(t)
			repo := repository.NewWeatherRepository(db)

			_, err := repo.GetWeather(tc.ctx, tc.id)
			require.ErrorAs(t, err, &repository.ErrNotFound{})
		})
	}
}

func TestUpdateWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		ctx       context.Context
		dbBuilder databaseBuilder
		ob        *models.Weather
		upd       *models.Weather
	}{
		{
			name: "Update third weather observation with extreme temperature",
			ctx:  context.Background(),
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(&models.Weather{
						ID:            1,
						Temperature:   -10.0,
						Humidity:      60.0,
						Pressure:      1025.0,
						WeatherStatus: "Blizzard",
						City:          "Helsinki",
						Country:       "Finland",
					}, nil).
					Once()
				mockService.
					On("UpdateWeather", mock.Anything, &models.Weather{
						ID:            1,
						Temperature:   -10.0,
						Humidity:      60.0,
						Pressure:      1025.0,
						WeatherStatus: "Blizzard",
						City:          "Helsinki",
						Country:       "Finland",
					}).
					Return(&models.Weather{
						ID:            1,
						Temperature:   -10.0,
						Humidity:      60.0,
						Pressure:      1025.0,
						WeatherStatus: "Blizzard",
						City:          "Helsinki",
						Country:       "Finland",
					}, nil).
					Once()
				mockService.
					On("GetWeather", mock.Anything, 1).
					Return(&models.Weather{
						ID:            1,
						Temperature:   -10.0,
						Humidity:      60.0,
						Pressure:      1025.0,
						WeatherStatus: "Blizzard",
						City:          "Helsinki",
						Country:       "Finland",
					}, nil).
					Once()

				return mockService
			},
			ob: &models.Weather{
				Temperature:   -5.0,
				Humidity:      55.0,
				Pressure:      1020.0,
				WeatherStatus: "Snowy",
				City:          "Helsinki",
				Country:       "Finland",
			},
			upd: &models.Weather{
				Temperature:   -10.0,
				Humidity:      60.0,
				Pressure:      1025.0,
				WeatherStatus: "Blizzard",
				City:          "Helsinki",
				Country:       "Finland",
			},
		},
		{
			name: "Update weather observation with new extreme humidity",
			ctx:  context.Background(),
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(&models.Weather{
						ID:            1,
						Temperature:   30.0,
						Humidity:      100.0,
						Pressure:      1015.0,
						WeatherStatus: "Torrential Rain",
						City:          "Singapore",
						Country:       "Singapore",
					}, nil).
					Once()
				mockService.
					On("UpdateWeather", mock.Anything, &models.Weather{
						ID:            1,
						Temperature:   30.0,
						Humidity:      100.0,
						Pressure:      1015.0,
						WeatherStatus: "Torrential Rain",
						City:          "Singapore",
						Country:       "Singapore",
					}).
					Return(&models.Weather{
						ID:            1,
						Temperature:   30.0,
						Humidity:      100.0,
						Pressure:      1015.0,
						WeatherStatus: "Torrential Rain",
						City:          "Singapore",
						Country:       "Singapore",
					}, nil).
					Once()
				mockService.
					On("GetWeather", mock.Anything, 1).
					Return(&models.Weather{
						ID:            1,
						Temperature:   30.0,
						Humidity:      100.0,
						Pressure:      1015.0,
						WeatherStatus: "Torrential Rain",
						City:          "Singapore",
						Country:       "Singapore",
					}, nil).
					Once()

				return mockService
			},
			ob: &models.Weather{
				Temperature:   28.0,
				Humidity:      95.0,
				Pressure:      1012.0,
				WeatherStatus: "Rainy",
				City:          "Singapore",
				Country:       "Singapore",
			},
			upd: &models.Weather{
				Temperature:   30.0,
				Humidity:      100.0,
				Pressure:      1015.0,
				WeatherStatus: "Torrential Rain",
				City:          "Singapore",
				Country:       "Singapore",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.dbBuilder(t)
			repo := repository.NewWeatherRepository(db)

			id, err := repo.AddWeather(tc.ctx, tc.ob)
			require.NoError(t, err)

			tc.upd.ID = id
			err = repo.UpdateWeather(tc.ctx, tc.upd)
			require.NoError(t, err)

			updatedOb, err := repo.GetWeather(tc.ctx, id)
			require.NoError(t, err)

			assert.Equal(t, tc.upd, updatedOb)
		})
	}
}

func TestUpdateWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		ctx       context.Context
		dbBuilder databaseBuilder
		id        int
	}{
		{
			name: "Update non-existing weather observation",
			ctx:  context.Background(),
			id:   10,
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("UpdateWeather", mock.Anything, &models.Weather{ID: 10}).
					Return(&models.Weather{}, repository.ErrNotFound{}).
					Once()

				return mockService
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.dbBuilder(t)
			repo := repository.NewWeatherRepository(db)

			err := repo.UpdateWeather(tc.ctx, &models.Weather{ID: tc.id})
			require.ErrorAs(t, err, &repository.ErrNotFound{})
		})
	}
}

func TestDeleteWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		ctx       context.Context
		dbBuilder databaseBuilder
		ob        *models.Weather
	}{
		{
			name: "Delete third weather observation with low temperature",
			ctx:  context.Background(),
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(&models.Weather{
						ID:            1,
						Temperature:   -5.0,
						Humidity:      55.0,
						Pressure:      1020.0,
						WeatherStatus: "Snowy",
						City:          "Helsinki",
						Country:       "Finland",
					}, nil).
					Once()
				mockService.
					On("DeleteWeather", mock.Anything, 1).
					Return(&models.Weather{
						ID:            1,
						Temperature:   -5.0,
						Humidity:      55.0,
						Pressure:      1020.0,
						WeatherStatus: "Snowy",
						City:          "Helsinki",
						Country:       "Finland",
					}, nil).
					Once()
				mockService.
					On("GetWeather", mock.Anything, 1).
					Return(nil, repository.ErrNotFound{}).
					Once()

				return mockService
			},
			ob: &models.Weather{
				Temperature:   -5.0,
				Humidity:      55.0,
				Pressure:      1020.0,
				WeatherStatus: "Snowy",
				City:          "Helsinki",
				Country:       "Finland",
			},
		},
		{
			name: "Delete weather observation with extreme humidity",
			ctx:  context.Background(),
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(&models.Weather{
						ID:            2,
						Temperature:   28.0,
						Humidity:      95.0,
						Pressure:      1012.0,
						WeatherStatus: "Rainy",
						City:          "Singapore",
						Country:       "Singapore",
					}, nil).
					Once()
				mockService.
					On("DeleteWeather", mock.Anything, 2).
					Return(&models.Weather{
						ID:            2,
						Temperature:   28.0,
						Humidity:      95.0,
						Pressure:      1012.0,
						WeatherStatus: "Rainy",
						City:          "Singapore",
						Country:       "Singapore",
					}, nil).
					Once()
				mockService.
					On("GetWeather", mock.Anything, 2).
					Return(nil, repository.ErrNotFound{}).
					Once()

				return mockService
			},
			ob: &models.Weather{
				Temperature:   28.0,
				Humidity:      95.0,
				Pressure:      1012.0,
				WeatherStatus: "Rainy",
				City:          "Singapore",
				Country:       "Singapore",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.dbBuilder(t)
			repo := repository.NewWeatherRepository(db)

			id, err := repo.AddWeather(tc.ctx, tc.ob)
			require.NoError(t, err)

			tc.ob.ID = id

			ob, err := repo.DeleteWeather(tc.ctx, id)
			require.NoError(t, err)
			assert.Equal(t, tc.ob, ob)

			_, err = repo.GetWeather(tc.ctx, id)
			assert.ErrorAs(t, err, &repository.ErrNotFound{})
		})
	}
}

func TestDeleteWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		ctx       context.Context
		dbBuilder databaseBuilder
		id        int
	}{
		{
			name: "Delete non-existent weather observation",
			ctx:  context.Background(),
			id:   1,
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("DeleteWeather", mock.Anything, 1).
					Return(nil, repository.ErrNotFound{}).
					Once()

				return mockService
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.dbBuilder(t)
			repo := repository.NewWeatherRepository(db)

			_, err := repo.DeleteWeather(tc.ctx, tc.id)
			require.ErrorAs(t, err, &repository.ErrNotFound{})
		})
	}
}

func TestListWeatherObservationsWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name         string
		dbBuilder    databaseBuilder
		observations []*models.Weather
	}{
		{
			name: "No weather observations",
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("ListWeathers", mock.Anything).
					Return([]*models.Weather{}, nil).
					Once()

				return mockService
			},
			observations: []*models.Weather{},
		},
		{
			name: "One weather observation",
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("ListWeathers", mock.Anything).
					Return([]*models.Weather{
						{
							ID:            1,
							Temperature:   25.0,
							Humidity:      60.0,
							Pressure:      1013.0,
							WeatherStatus: "Clear",
							City:          "Berlin",
							Country:       "Germany",
						},
					}, nil).
					Once()

				return mockService
			},
			observations: []*models.Weather{
				{
					ID:            1,
					Temperature:   25.0,
					Humidity:      60.0,
					Pressure:      1013.0,
					WeatherStatus: "Clear",
					City:          "Berlin",
					Country:       "Germany",
				},
			},
		},
		{
			name: "Multiple weather observations",
			dbBuilder: func(t *testing.T) repository.Database {
				t.Helper()

				mockService := NewMockDatabase(t)
				mockService.
					On("ListWeathers", mock.Anything).
					Return([]*models.Weather{
						{
							ID:            1,
							Temperature:   25.0,
							Humidity:      60.0,
							Pressure:      1013.0,
							WeatherStatus: "Clear",
							City:          "Berlin",
							Country:       "Germany",
						},
						{
							ID:            2,
							Temperature:   30.0,
							Humidity:      65.0,
							Pressure:      1015.0,
							WeatherStatus: "Sunny",
							City:          "Paris",
							Country:       "France",
						},
						{
							ID:            3,
							Temperature:   28.0,
							Humidity:      95.0,
							Pressure:      1012.0,
							WeatherStatus: "Rainy",
							City:          "Singapore",
							Country:       "Singapore",
						},
					}, nil).
					Once()

				return mockService
			},
			observations: []*models.Weather{
				{
					ID:            1,
					Temperature:   25.0,
					Humidity:      60.0,
					Pressure:      1013.0,
					WeatherStatus: "Clear",
					City:          "Berlin",
					Country:       "Germany",
				},
				{
					ID:            2,
					Temperature:   30.0,
					Humidity:      65.0,
					Pressure:      1015.0,
					WeatherStatus: "Sunny",
					City:          "Paris",
					Country:       "France",
				},
				{
					ID:            3,
					Temperature:   28.0,
					Humidity:      95.0,
					Pressure:      1012.0,
					WeatherStatus: "Rainy",
					City:          "Singapore",
					Country:       "Singapore",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.dbBuilder(t)
			repo := repository.NewWeatherRepository(db)

			observations, err := repo.ListWeathers(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.observations, observations)
		})
	}
}
