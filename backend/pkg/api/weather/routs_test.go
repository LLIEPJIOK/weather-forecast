package weather_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/backend/internal/repository"
	"github.com/LLIEPJIOK/weather-forecast/backend/pkg/api/weather"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type serviceBuilder func(t *testing.T) weather.WeatherService

func TestAddWeatherHandlerWithBuilder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		inputBody          string
		repoBuilder        serviceBuilder
		expectedStatusCode int
		expectedResponse   string
	}

	tt := []testCase{
		{
			name:      "Valid input",
			inputBody: `{"Temperature": 25.5, "Humidity": 80}`,
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(1, nil).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":1}`,
		},
		{
			name:      "Invalid JSON",
			inputBody: `{"invalid_json":}`,
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				return NewMockWeatherService(t)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"invalid input: code=400, message=Syntax error: offset=17, error=invalid character '}' looking for beginning of value, internal=invalid character '}' looking for beginning of value"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"Temperature": 123, "Humidity": 10}`,
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("AddWeather", mock.Anything, mock.Anything).
					Return(0, errors.New("database error")).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"The server is temporarily unavailable, please try again later"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := tc.repoBuilder(t)

			e := echo.New()

			req := httptest.NewRequest(
				http.MethodPost,
				"/weather",
				bytes.NewReader([]byte(tc.inputBody)),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler := weather.AddWeatherHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestGetWeatherHandlerWithBuilder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		inputID            string
		repoBuilder        serviceBuilder
		expectedStatusCode int
		expectedResponse   string
	}

	tm := time.Now()

	tt := []testCase{
		{
			name:    "Valid ID",
			inputID: "1",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("GetWeather", mock.Anything, 1).
					Return(&models.Weather{
						ID:            1,
						City:          "Berlin",
						Country:       "Germany",
						Timestamp:     tm,
						Temperature:   25.5,
						Humidity:      80,
						Pressure:      1013,
						WindSpeed:     5.4,
						WeatherStatus: "Clear",
					}, nil).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
                "id": 1,
								"city": "Berlin",
								"country": "Germany",
                "timestamp": "` + tm.Format(time.RFC3339Nano) + `",
                "temperature": 25.5,
                "humidity": 80,
                "pressure": 1013,
								"wind_speed": 5.4,
                "weather_status": "Clear"
            }`,
		},
		{
			name:    "Invalid ID",
			inputID: "invalid_id",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()
				return NewMockWeatherService(t)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
				"message": "parseID: failed to parse id=\"invalid_id\": strconv.Atoi: parsing \"invalid_id\": invalid syntax"
			}`,
		},
		{
			name:    "Not Found",
			inputID: "2",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("GetWeather", mock.Anything, 2).
					Return(&models.Weather{}, repository.ErrNotFound{}).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: `{
				"message": "record with id=2 not found"
			}`,
		},
		{
			name:    "Service error",
			inputID: "3",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("GetWeather", mock.Anything, 3).
					Return(&models.Weather{}, errors.New("database error")).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: `{
				"message": "The server is temporarily unavailable, please try again later"
			}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := tc.repoBuilder(t)

			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/weather", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/weather/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.inputID)

			handler := weather.GetWeatherHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestUpdateWeatherHandlerWithBuilder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		inputID            string
		inputWeatherObs    *models.Weather
		repoBuilder        serviceBuilder
		expectedStatusCode int
		expectedResponse   string
	}

	tt := []testCase{
		{
			name:    "Valid ID and valid data",
			inputID: "1",
			inputWeatherObs: &models.Weather{
				Temperature:   22.5,
				Humidity:      70,
				Pressure:      1012,
				WindSpeed:     3.5,
				WeatherStatus: "Cloudy",
			},
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On(
						"UpdateWeather",
						mock.Anything,
						&models.Weather{
							ID:            1,
							Temperature:   22.5,
							Humidity:      70,
							Pressure:      1012,
							WindSpeed:     3.5,
							WeatherStatus: "Cloudy",
						},
					).
					Return(nil).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
				"message": "successfully updated"
			}`,
		},
		{
			name:            "Invalid ID",
			inputID:         "invalid_id",
			inputWeatherObs: &models.Weather{},
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()
				return NewMockWeatherService(t)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
				"message": "parseID: failed to parse id=\"invalid_id\": strconv.Atoi: parsing \"invalid_id\": invalid syntax"
			}`,
		},
		{
			name:    "Not Found",
			inputID: "2",
			inputWeatherObs: &models.Weather{
				Temperature:   22.5,
				Humidity:      70,
				Pressure:      1012,
				WindSpeed:     3.5,
				WeatherStatus: "Cloudy",
			},
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On(
						"UpdateWeather",
						mock.Anything,
						&models.Weather{
							ID:            2,
							Temperature:   22.5,
							Humidity:      70,
							Pressure:      1012,
							WindSpeed:     3.5,
							WeatherStatus: "Cloudy",
						},
					).
					Return(repository.ErrNotFound{}).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: `{
				"message": "record with id=2 not found"
			}`,
		},
		{
			name:    "Service error",
			inputID: "3",
			inputWeatherObs: &models.Weather{
				Temperature:   22.5,
				Humidity:      70,
				Pressure:      1012,
				WindSpeed:     3.5,
				WeatherStatus: "Cloudy",
			},
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On(
						"UpdateWeather",
						mock.Anything,
						&models.Weather{
							ID:            3,
							Temperature:   22.5,
							Humidity:      70,
							Pressure:      1012,
							WindSpeed:     3.5,
							WeatherStatus: "Cloudy",
						},
					).
					Return(errors.New("database error")).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: `{
				"message": "The server is temporarily unavailable, please try again later"
			}`,
		},
		{
			name:            "Invalid JSON format",
			inputID:         "1",
			inputWeatherObs: &models.Weather{},
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()
				return NewMockWeatherService(t)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
				"message": "invalid input: code=400, message=Unmarshal type error: expected=float64, got=string, field=temperature, offset=31, internal=json: cannot unmarshal string into Go struct field Weather.temperature of type float64"
			}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := tc.repoBuilder(t)

			e := echo.New()

			var reqBody []byte
			if tc.name == "Invalid JSON format" {
				reqBody = []byte(
					`{"Temperature": "invalid_value", "Humidity": 70}`,
				)
			} else {
				var err error

				reqBody, err = json.Marshal(*tc.inputWeatherObs)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/weather", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/weather/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.inputID)

			handler := weather.UpdateWeatherHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestDeleteWeatherHandlerWithBuilder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		inputID            string
		repoBuilder        serviceBuilder
		expectedStatusCode int
		expectedResponse   string
	}

	tm := time.Now()

	tt := []testCase{
		{
			name:    "Valid ID",
			inputID: "1",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("DeleteWeather", mock.Anything, 1).
					Return(&models.Weather{
						ID:            1,
						City:          "Berlin",
						Country:       "Germany",
						Timestamp:     tm,
						Temperature:   25.5,
						Humidity:      80,
						Pressure:      1013,
						WindSpeed:     5.4,
						WeatherStatus: "Clear",
					}, nil).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
				"id": 1,
				"city": "Berlin",
				"country": "Germany",
				"timestamp": "` + tm.Format(time.RFC3339Nano) + `",
				"temperature": 25.5,
				"humidity": 80,
				"pressure": 1013,
				"wind_speed": 5.4,
				"weather_status": "Clear"
			}`,
		},
		{
			name:    "Invalid ID",
			inputID: "invalid_id",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()
				return NewMockWeatherService(t)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
				"message": "parseID: failed to parse id=\"invalid_id\": strconv.Atoi: parsing \"invalid_id\": invalid syntax"
			}`,
		},
		{
			name:    "Not Found",
			inputID: "2",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("DeleteWeather", mock.Anything, 2).
					Return(&models.Weather{}, repository.ErrNotFound{}).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: `{
				"message": "record with id=2 not found"
			}`,
		},
		{
			name:    "Service error",
			inputID: "3",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("DeleteWeather", mock.Anything, 3).
					Return(&models.Weather{}, errors.New("database error")).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: `{
				"message": "The server is temporarily unavailable, please try again later"
			}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := tc.repoBuilder(t)

			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/weather", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/weather/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.inputID)

			handler := weather.DeleteWeatherHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestListWeathersHandlerWithBuilder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		repoBuilder        serviceBuilder
		expectedStatusCode int
		expectedResponse   string
	}

	tm := time.Now()

	tt := []testCase{
		{
			name: "Valid weathers",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("ListWeathers", mock.Anything).
					Return([]*models.Weather{
						{
							ID:            1,
							City:          "Berlin",
							Country:       "Germany",
							Timestamp:     tm,
							Temperature:   25.5,
							Humidity:      80,
							Pressure:      1013,
							WindSpeed:     5.4,
							WeatherStatus: "Clear",
						},
						{
							ID:            2,
							City:          "Paris",
							Country:       "France",
							Timestamp:     tm.Add(-time.Hour * 24 * 365),
							Temperature:   15.0,
							Humidity:      75,
							Pressure:      1012,
							WindSpeed:     3.2,
							WeatherStatus: "Cloudy",
						},
					}, nil).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `[
				{
					"id": 1,
					"city": "Berlin",
					"country": "Germany",
					"timestamp": "` + tm.Format(time.RFC3339Nano) + `",
					"temperature": 25.5,
					"humidity": 80,
					"pressure": 1013,
					"wind_speed": 5.4,
					"weather_status": "Clear"
				},
				{
					"id": 2,
					"city": "Paris",
					"country": "France",
					"timestamp": "` + tm.Add(-time.Hour*24*365).Format(time.RFC3339Nano) + `",
					"temperature": 15.0,
					"humidity": 75,
					"pressure": 1012,
					"wind_speed": 3.2,
					"weather_status": "Cloudy"
				}
			]`,
		},
		{
			name: "Service error",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("ListWeathers", mock.Anything).
					Return(nil, errors.New("database error")).
					Once()

				return mockService
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: `{
				"message": "The server is temporarily unavailable, please try again later"
			}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := tc.repoBuilder(t)

			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/weathers", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			handler := weather.ListWeathersHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
