package repository_test

import (
	"context"
	"testing"

	"github.com/LLIEPJIOK/weather-forecast/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		ctx  context.Context
		ob   models.WeatherObservation
	}{
		{
			name: "Add first weather observation with positive values",
			ctx:  context.Background(),
			ob: models.WeatherObservation{
				Temperature:   25.0,
				Humidity:      60.0,
				Pressure:      1013.0,
				WeatherStatus: "Clear",
				Location: models.Location{
					Latitude:  52.5200,
					Longitude: 13.4050,
					City:      "Berlin",
					Country:   "Germany",
				},
			},
		},
		{
			name: "Add second weather observation with different location",
			ctx:  context.Background(),
			ob: models.WeatherObservation{
				Temperature:   30.0,
				Humidity:      65.0,
				Pressure:      1015.0,
				WeatherStatus: "Sunny",
				Location: models.Location{
					Latitude:  48.8566,
					Longitude: 2.3522,
					City:      "Paris",
					Country:   "France",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewWeatherRepository()

			id, err := repo.AddWeatherObservation(tc.ctx, tc.ob)
			require.NoError(t, err)

			tc.ob.ID = id

			ob, err := repo.GetWeatherObservation(tc.ctx, id)
			require.NoError(t, err)

			assert.Equal(t, tc.ob, ob)
		})
	}
}

func TestGetWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		ctx  context.Context
		id   int
	}{
		{
			name: "Update third weather observation with extreme temperature",
			ctx:  context.Background(),
			id:   1234,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewWeatherRepository()

			_, err := repo.GetWeatherObservation(tc.ctx, tc.id)
			require.ErrorAs(t, err, &repository.ErrNotFound{})
		})
	}
}

func TestUpdateWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		ctx  context.Context
		ob   models.WeatherObservation
		upd  models.WeatherObservation
	}{
		{
			name: "Update third weather observation with extreme temperature",
			ctx:  context.Background(),
			ob: models.WeatherObservation{
				Temperature:   -5.0,
				Humidity:      55.0,
				Pressure:      1020.0,
				WeatherStatus: "Snowy",
				Location: models.Location{
					Latitude:  60.1692,
					Longitude: 24.9402,
					City:      "Helsinki",
					Country:   "Finland",
				},
			},
			upd: models.WeatherObservation{
				Temperature:   -10.0,
				Humidity:      60.0,
				Pressure:      1025.0,
				WeatherStatus: "Blizzard",
				Location: models.Location{
					Latitude:  60.1692,
					Longitude: 24.9402,
					City:      "Helsinki",
					Country:   "Finland",
				},
			},
		},
		{
			name: "Update weather observation with new extreme humidity",
			ctx:  context.Background(),
			ob: models.WeatherObservation{
				Temperature:   28.0,
				Humidity:      95.0,
				Pressure:      1012.0,
				WeatherStatus: "Rainy",
				Location: models.Location{
					Latitude:  1.3521,
					Longitude: 103.8198,
					City:      "Singapore",
					Country:   "Singapore",
				},
			},
			upd: models.WeatherObservation{
				Temperature:   30.0,
				Humidity:      100.0,
				Pressure:      1015.0,
				WeatherStatus: "Torrential Rain",
				Location: models.Location{
					Latitude:  1.3521,
					Longitude: 103.8198,
					City:      "Singapore",
					Country:   "Singapore",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewWeatherRepository()

			id, err := repo.AddWeatherObservation(tc.ctx, tc.ob)
			require.NoError(t, err)

			tc.upd.ID = id
			err = repo.UpdateWeatherObservation(tc.ctx, tc.upd)
			require.NoError(t, err)

			updatedOb, err := repo.GetWeatherObservation(tc.ctx, id)
			require.NoError(t, err)

			assert.Equal(t, tc.upd, updatedOb)
		})
	}
}

func TestUpdateWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		ctx  context.Context
		id   int
	}{
		{
			name: "Update third weather observation with extreme temperature",
			ctx:  context.Background(),
			id:   10,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewWeatherRepository()

			err := repo.UpdateWeatherObservation(tc.ctx, models.WeatherObservation{ID: tc.id})
			require.ErrorAs(t, err, &repository.ErrNotFound{})
		})
	}
}

func TestDeleteWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		ctx  context.Context
		ob   models.WeatherObservation
	}{
		{
			name: "Delete third weather observation with low temperature",
			ctx:  context.Background(),
			ob: models.WeatherObservation{
				Temperature:   -5.0,
				Humidity:      55.0,
				Pressure:      1020.0,
				WeatherStatus: "Snowy",
				Location: models.Location{
					Latitude:  60.1692,
					Longitude: 24.9402,
					City:      "Helsinki",
					Country:   "Finland",
				},
			},
		},
		{
			name: "Delete weather observation with extreme humidity",
			ctx:  context.Background(),
			ob: models.WeatherObservation{
				Temperature:   28.0,
				Humidity:      95.0,
				Pressure:      1012.0,
				WeatherStatus: "Rainy",
				Location: models.Location{
					Latitude:  1.3521,
					Longitude: 103.8198,
					City:      "Singapore",
					Country:   "Singapore",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewWeatherRepository()

			id, err := repo.AddWeatherObservation(tc.ctx, tc.ob)
			require.NoError(t, err)

			tc.ob.ID = id

			ob, err := repo.DeleteWeatherObservation(tc.ctx, id)
			require.NoError(t, err)
			assert.Equal(t, tc.ob, ob)

			_, err = repo.GetWeatherObservation(tc.ctx, id)
			assert.Error(t, err)
		})
	}
}

func TestDeleteWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		ctx  context.Context
		id   int
	}{
		{
			name: "Delete non-existent weather observation",
			ctx:  context.Background(),
			id:   1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewWeatherRepository()

			_, err := repo.DeleteWeatherObservation(tc.ctx, tc.id)
			require.ErrorAs(t, err, &repository.ErrNotFound{})
		})
	}
}

func TestListWeatherObservationsWithoutError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name         string
		observations []models.WeatherObservation
	}{
		{
			name:         "No weather observations",
			observations: nil,
		},
		{
			name: "One weather observation",
			observations: []models.WeatherObservation{
				{
					Temperature:   25.0,
					Humidity:      60.0,
					Pressure:      1013.0,
					WeatherStatus: "Clear",
					Location: models.Location{
						Latitude:  52.5200,
						Longitude: 13.4050,
						City:      "Berlin",
						Country:   "Germany",
					},
				},
			},
		},
		{
			name: "Multiple weather observations",
			observations: []models.WeatherObservation{
				{
					Temperature:   25.0,
					Humidity:      60.0,
					Pressure:      1013.0,
					WeatherStatus: "Clear",
					Location: models.Location{
						Latitude:  52.5200,
						Longitude: 13.4050,
						City:      "Berlin",
						Country:   "Germany",
					},
				},
				{
					Temperature:   30.0,
					Humidity:      65.0,
					Pressure:      1015.0,
					WeatherStatus: "Sunny",
					Location: models.Location{
						Latitude:  48.8566,
						Longitude: 2.3522,
						City:      "Paris",
						Country:   "France",
					},
				},
				{
					Temperature:   28.0,
					Humidity:      95.0,
					Pressure:      1012.0,
					WeatherStatus: "Rainy",
					Location: models.Location{
						Latitude:  1.3521,
						Longitude: 103.8198,
						City:      "Singapore",
						Country:   "Singapore",
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewWeatherRepository()

			for i, ob := range tc.observations {
				id, err := repo.AddWeatherObservation(context.Background(), ob)
				require.NoError(t, err)

				tc.observations[i].ID = id
			}

			observations, err := repo.ListWeatherObservations(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.observations, observations)
		})
	}
}
