package weather_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LLIEPJIOK/weather-forecast/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/internal/repository"
	"github.com/LLIEPJIOK/weather-forecast/pkg/api/weather"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type serviceBuilder func(t *testing.T) weather.WeatherService

func TestAddWeatherObservationHandlerWithBuilder(t *testing.T) {
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
					On("AddWeatherObservation", mock.Anything, mock.Anything).
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
					On("AddWeatherObservation", mock.Anything, mock.Anything).
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
			handler := weather.AddWeatherObservationHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestGetWeatherObservationHandlerWithBuilder(t *testing.T) {
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
					On("GetWeatherObservation", mock.Anything, 1).
					Return(models.WeatherObservation{
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
					On("GetWeatherObservation", mock.Anything, 2).
					Return(models.WeatherObservation{}, repository.ErrNotFound{}).
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
					On("GetWeatherObservation", mock.Anything, 3).
					Return(models.WeatherObservation{}, errors.New("database error")).
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

			handler := weather.GetWeatherObservationHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestUpdateWeatherObservationHandlerWithBuilder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		inputID            string
		inputWeatherObs    models.WeatherObservation
		repoBuilder        serviceBuilder
		expectedStatusCode int
		expectedResponse   string
	}

	tt := []testCase{
		{
			name:    "Valid ID and valid data",
			inputID: "1",
			inputWeatherObs: models.WeatherObservation{
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
						"UpdateWeatherObservation",
						mock.Anything,
						mock.MatchedBy(func(obs models.WeatherObservation) bool {
							return obs.ID == 1 && obs.Temperature == 22.5 && obs.Humidity == 70 &&
								obs.Pressure == 1012
						}),
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
			inputWeatherObs: models.WeatherObservation{
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
						"UpdateWeatherObservation",
						mock.Anything,
						mock.MatchedBy(func(obs models.WeatherObservation) bool {
							return obs.ID == 2
						}),
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
			inputWeatherObs: models.WeatherObservation{
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
						"UpdateWeatherObservation",
						mock.Anything,
						mock.MatchedBy(func(obs models.WeatherObservation) bool {
							return obs.ID == 3
						}),
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
			inputWeatherObs: models.WeatherObservation{},
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()
				return NewMockWeatherService(t)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
				"message": "invalid input: code=400, message=Unmarshal type error: expected=float64, got=string, field=temperature, offset=31, internal=json: cannot unmarshal string into Go struct field WeatherObservation.temperature of type float64"
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
				) // Невалидный JSON
			} else {
				var err error
				reqBody, err = json.Marshal(tc.inputWeatherObs)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/weather", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/weather/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.inputID)

			handler := weather.UpdateWeatherObservationHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestDeleteWeatherObservationHandlerWithBuilder(t *testing.T) {
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
					On("DeleteWeatherObservation", mock.Anything, 1).
					Return(models.WeatherObservation{
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
					On("DeleteWeatherObservation", mock.Anything, 2).
					Return(models.WeatherObservation{}, repository.ErrNotFound{}).
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
					On("DeleteWeatherObservation", mock.Anything, 3).
					Return(models.WeatherObservation{}, errors.New("database error")).
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

			handler := weather.DeleteWeatherObservationHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func TestListWeatherObservationsHandlerWithBuilder(t *testing.T) {
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
			name: "Valid observations",
			repoBuilder: func(t *testing.T) weather.WeatherService {
				t.Helper()

				mockService := NewMockWeatherService(t)
				mockService.
					On("ListWeatherObservations", mock.Anything).
					Return([]models.WeatherObservation{
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
					On("ListWeatherObservations", mock.Anything).
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

			handler := weather.ListWeatherObservationsHandler(mockService)

			err := handler(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
