package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddWeatherWithoutError(t *testing.T) {
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
				repo.On("AddWeather", mock.Anything, mock.Anything).Return(1, nil).Once()

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

			id, err := srv.AddWeather(tc.ctx, &models.Weather{})
			require.NoError(t, err)
			assert.Equal(t, tc.id, id)
		})
	}
}

func TestAddWeatherWithError(t *testing.T) {
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
				repo.On("AddWeather", mock.Anything, mock.Anything).
					Return(0, errAdd).
					Once()

				return repo
			},
			ctx: context.Background(),
			err: fmt.Errorf("failed to add weather: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.AddWeather(tc.ctx, &models.Weather{})
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestGetWeatherWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		id          int
		ob          *models.Weather
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("GetWeather", mock.Anything, 1).Return(&models.Weather{
					Temperature:   25.0,
					Humidity:      60.0,
					Pressure:      1013.0,
					WeatherStatus: "Clear",
					City:          "Berlin",
					Country:       "Germany",
				}, nil).Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
			ob: &models.Weather{
				Temperature:   25.0,
				Humidity:      60.0,
				Pressure:      1013.0,
				WeatherStatus: "Clear",
				City:          "Berlin",
				Country:       "Germany",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			ob, err := srv.GetWeather(tc.ctx, tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.ob, ob)
		})
	}
}

func TestGetWeatherWithError(t *testing.T) {
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
				repo.On("GetWeather", mock.Anything, 1).
					Return(&models.Weather{}, errGet).
					Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
			err: fmt.Errorf("failed to get weather: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.GetWeather(tc.ctx, tc.id)
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestUpdateWeatherWithoutError(t *testing.T) {
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
				repo.On("UpdateWeather", mock.Anything, mock.Anything).Return(nil).Once()

				return repo
			},
			ctx: context.Background(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			err := srv.UpdateWeather(tc.ctx, &models.Weather{})
			require.NoError(t, err)
		})
	}
}

func TestUpdateWeatherWithError(t *testing.T) {
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
				repo.On("UpdateWeather", mock.Anything, mock.Anything).
					Return(errUpdate).
					Once()

				return repo
			},
			ctx: context.Background(),
			err: fmt.Errorf("failed to update weather: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			err := srv.UpdateWeather(tc.ctx, &models.Weather{})
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestDeleteWeatherWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		id          int
		ob          *models.Weather
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("DeleteWeather", mock.Anything, 1).
					Return(&models.Weather{
						Temperature:   30.0,
						Humidity:      65.0,
						Pressure:      1015.0,
						WeatherStatus: "Sunny",
						City:          "Paris",
						Country:       "France",
					}, nil).
					Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
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

			srv := service.NewWeatherService(tc.repoBuilder(t))

			ob, err := srv.DeleteWeather(tc.ctx, tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.ob, ob)
		})
	}
}

func TestDeleteWeatherWithError(t *testing.T) {
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
				repo.On("DeleteWeather", mock.Anything, 1).
					Return(&models.Weather{}, errDelete).
					Once()

				return repo
			},
			ctx: context.Background(),
			id:  1,
			err: fmt.Errorf("failed to delete weather: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.DeleteWeather(tc.ctx, tc.id)
			require.EqualError(t, err, tc.err.Error())
		})
	}
}

func TestListWeathersWithoutError(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		repoBuilder func(t *testing.T) service.WeatherRepo
		ctx         context.Context
		obList      []*models.Weather
	}

	tt := []TestCase{
		{
			name: "primary",
			repoBuilder: func(t *testing.T) service.WeatherRepo {
				t.Helper()

				repo := NewMockWeatherRepo(t)
				repo.On("ListWeathers", mock.Anything).
					Return([]*models.Weather{
						{
							Temperature:   25.0,
							Humidity:      60.0,
							Pressure:      1013.0,
							WeatherStatus: "Clear",
							City:          "Berlin",
							Country:       "Germany",
						},
						{
							Temperature:   30.0,
							Humidity:      65.0,
							Pressure:      1015.0,
							WeatherStatus: "Sunny",
							City:          "Paris",
							Country:       "France",
						},
						{
							Temperature:   28.0,
							Humidity:      95.0,
							Pressure:      1012.0,
							WeatherStatus: "Rainy",
							City:          "Singapore",
							Country:       "Singapore",
						},
					}, nil).
					Once()

				return repo
			},
			ctx: context.Background(),
			obList: []*models.Weather{
				{
					Temperature:   25.0,
					Humidity:      60.0,
					Pressure:      1013.0,
					WeatherStatus: "Clear",
					City:          "Berlin",
					Country:       "Germany",
				},
				{
					Temperature:   30.0,
					Humidity:      65.0,
					Pressure:      1015.0,
					WeatherStatus: "Sunny",
					City:          "Paris",
					Country:       "France",
				},
				{
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

			srv := service.NewWeatherService(tc.repoBuilder(t))

			obList, err := srv.ListWeathers(tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.obList, obList)
		})
	}
}

func TestListWeathersWithError(t *testing.T) {
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
				repo.On("ListWeathers", mock.Anything).Return(nil, errList).Once()

				return repo
			},
			ctx: context.Background(),
			err: fmt.Errorf("failed to list weathers: repo error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := service.NewWeatherService(tc.repoBuilder(t))

			_, err := srv.ListWeathers(tc.ctx)
			require.EqualError(t, err, tc.err.Error())
		})
	}
}
