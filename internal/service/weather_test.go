package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/LLIEPJIOK/weather-forecast/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		id          int
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("AddWeatherObservation", mock.Anything, mock.Anything).Return(1, nil).Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			id, err := srv.AddWeatherObservation(tc.ctx, models.WeatherObservation{})
			require.NoError(t, err)
			assert.Equal(t, tc.id, id)
		})
	}
}

func TestAddWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		err         error
	}

	errAdd := fmt.Errorf("repo error")

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("AddWeatherObservation", mock.Anything, mock.Anything).
					Return(0, errAdd).
					Once()

				return repo
			},
			ctx: context.Background(),
			err: fmt.Errorf("failed to add weather observation: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.AddWeatherObservation(tc.ctx, models.WeatherObservation{})
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestGetWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		id          int
		ob          models.WeatherObservation
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("GetWeatherObservation", mock.Anything, 1).Return(models.WeatherObservation{
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
				}, nil).Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			ob, err := srv.GetWeatherObservation(tc.ctx, tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.ob, ob)
		})
	}
}

func TestGetWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		id          int
		err         error
	}

	errGet := fmt.Errorf("repo error")

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("GetWeatherObservation", mock.Anything, 1).
					Return(models.WeatherObservation{}, errGet).
					Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
			err: fmt.Errorf("failed to get weather observation: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.GetWeatherObservation(tc.ctx, tc.id)
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestUpdateWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("UpdateWeatherObservation", mock.Anything, mock.Anything).Return(nil).Once()

				return repo
			},
			ctx: context.Background(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			err := srv.UpdateWeatherObservation(tc.ctx, models.WeatherObservation{})
			require.NoError(t, err)
		})
	}
}

func TestUpdateWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		err         error
	}

	errUpdate := fmt.Errorf("repo error")

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("UpdateWeatherObservation", mock.Anything, mock.Anything).
					Return(errUpdate).
					Once()

				return repo
			},
			ctx: context.Background(),
			err: fmt.Errorf("failed to update weather observation: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			err := srv.UpdateWeatherObservation(tc.ctx, models.WeatherObservation{})
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestDeleteWeatherObservationWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		id          int
		ob          models.WeatherObservation
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("DeleteWeatherObservation", mock.Anything, 1).
					Return(models.WeatherObservation{
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
					}, nil).
					Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
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

			srv := service.NewWeatherService(tc.repoBuilder(t))

			ob, err := srv.DeleteWeatherObservation(tc.ctx, tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.ob, ob)
		})
	}
}

func TestDeleteWeatherObservationWithError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		id          int
		err         error
	}

	errDelete := fmt.Errorf("repo error")

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("DeleteWeatherObservation", mock.Anything, 1).
					Return(models.WeatherObservation{}, errDelete).
					Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
			err: fmt.Errorf("failed to delete weather observation: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.DeleteWeatherObservation(tc.ctx, tc.id)
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestListWeatherObservationsWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		obList      []models.WeatherObservation
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("ListWeatherObservations", mock.Anything).
					Return([]models.WeatherObservation{
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
					}, nil).
					Once()

				return repo
			},
			ctx: context.Background(),
			obList: []models.WeatherObservation{
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

			srv := service.NewWeatherService(tc.repoBuilder(t))

			obList, err := srv.ListWeatherObservations(tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.obList, obList)
		})
	}
}

func TestListWeatherObservationsWithError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		err         error
	}

	errList := fmt.Errorf("repo error")

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("ListWeatherObservations", mock.Anything).Return(nil, errList).Once()

				return repo
			},
			ctx: context.Background(),
			err: fmt.Errorf("failed to list weather observations: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.ListWeatherObservations(tc.ctx)
			require.EqualError(t, err, tc.err.Error())
		})
	}
}
