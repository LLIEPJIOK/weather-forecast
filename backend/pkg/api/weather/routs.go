package weather

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/backend/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoMessage struct {
	Msg string `json:"message"`
}

type EchoID struct {
	ID int `json:"id"`
}

//go:generate mockery --name WeatherService --structname MockWeatherService --filename mock_weather_service_test.go --outpkg weather_test --output .
type WeatherService interface {
	AddWeather(ctx context.Context, ob *models.Weather) (int, error)
	GetWeather(ctx context.Context, id int) (*models.Weather, error)
	UpdateWeather(ctx context.Context, ob *models.Weather) error
	DeleteWeather(ctx context.Context, id int) (*models.Weather, error)
	ListWeathers(ctx context.Context) ([]*models.Weather, error)
}

func RegisterWeatherRoutes(
	ctx context.Context,
	server *echo.Echo,
	weatherService WeatherService,
) {
	server.POST("/weather", AddWeatherHandler(weatherService))
	server.GET("/weather/:id", GetWeatherHandler(weatherService))
	server.PUT("/weather/:id", UpdateWeatherHandler(weatherService))
	server.DELETE("/weather/:id", DeleteWeatherHandler(weatherService))
	server.GET("/weathers", ListWeathersHandler(weatherService))

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // Разрешаем запросы с вашего фронтенда
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
}

func AddWeatherHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ob models.Weather

		if err := c.Bind(&ob); err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("invalid input: %s", err)},
				"\t",
			)
		}

		id, err := weatherService.AddWeather(c.Request().Context(), &ob)
		if err != nil {
			c.Logger().Errorf("failed to add: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, EchoID{ID: id}, "\t")
	}
}

func GetWeatherHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseID(c)
		if err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("parseID: %s", err)},
				"\t",
			)
		}

		ob, err := weatherService.GetWeather(c.Request().Context(), id)
		if err != nil {
			if errors.As(err, &repository.ErrNotFound{}) {
				return c.JSONPretty(
					http.StatusNotFound,
					EchoMessage{Msg: fmt.Sprintf("record with id=%d not found", id)},
					"\t",
				)
			}

			c.Logger().Errorf("failed to get: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, ob, "\t")
	}
}

func UpdateWeatherHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ob models.Weather

		if err := c.Bind(&ob); err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("invalid input: %s", err)},
				"\t",
			)
		}

		id, err := parseID(c)
		if err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("parseID: %s", err)},
				"\t",
			)
		}

		ob.ID = id

		err = weatherService.UpdateWeather(c.Request().Context(), &ob)
		if err != nil {
			if errors.As(err, &repository.ErrNotFound{}) {
				return c.JSONPretty(
					http.StatusNotFound,
					EchoMessage{Msg: fmt.Sprintf("record with id=%d not found", id)},
					"\t",
				)
			}

			c.Logger().Errorf("failed to update: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, EchoMessage{Msg: "successfully updated"}, "\t")
	}
}

func DeleteWeatherHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseID(c)
		if err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("parseID: %s", err)},
				"\t",
			)
		}

		ob, err := weatherService.DeleteWeather(c.Request().Context(), id)
		if err != nil {
			if errors.As(err, &repository.ErrNotFound{}) {
				return c.JSONPretty(
					http.StatusNotFound,
					EchoMessage{Msg: fmt.Sprintf("record with id=%d not found", id)},
					"\t",
				)
			}

			c.Logger().Errorf("failed to delete: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, ob, "\t")
	}
}

func ListWeathersHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		obs, err := weatherService.ListWeathers(c.Request().Context())
		if err != nil {
			c.Logger().Errorf("failed to get list of weathers: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, obs, "\t")
	}
}

func parseID(c echo.Context) (int, error) {
	strID := c.Param("id")

	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, fmt.Errorf("failed to parse id=%q: %w", strID, err)
	}

	return id, nil
}
