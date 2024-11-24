package weather

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LLIEPJIOK/weather-forecast/internal/models"
	"github.com/LLIEPJIOK/weather-forecast/internal/repository"

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
	AddWeatherObservation(ctx context.Context, ob models.WeatherObservation) (int, error)
	GetWeatherObservation(ctx context.Context, id int) (models.WeatherObservation, error)
	UpdateWeatherObservation(ctx context.Context, ob models.WeatherObservation) error
	DeleteWeatherObservation(ctx context.Context, id int) (models.WeatherObservation, error)
	ListWeatherObservations(ctx context.Context) ([]models.WeatherObservation, error)
}

func RegisterWeatherObservationRoutes(
	ctx context.Context,
	server *echo.Echo,
	weatherService WeatherService,
) {
	server.POST("/weather", AddWeatherObservationHandler(weatherService))
	server.GET("/weather/:id", GetWeatherObservationHandler(weatherService))
	server.PUT("/weather/:id", UpdateWeatherObservationHandler(weatherService))
	server.DELETE("/weather/:id", DeleteWeatherObservationHandler(weatherService))
	server.GET("/weathers", ListWeatherObservationsHandler(weatherService))

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // Разрешаем запросы с вашего фронтенда
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
}

func AddWeatherObservationHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ob models.WeatherObservation

		if err := c.Bind(&ob); err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("invalid input: %s", err)},
				"\t",
			)
		}

		id, err := weatherService.AddWeatherObservation(c.Request().Context(), ob)
		if err != nil {
			c.Logger().Errorf("failed to add observation: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, EchoID{ID: id}, "\t")
	}
}

func GetWeatherObservationHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseID(c)
		if err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("parseID: %s", err)},
				"\t",
			)
		}

		ob, err := weatherService.GetWeatherObservation(c.Request().Context(), id)
		if err != nil {
			if errors.As(err, &repository.ErrNotFound{}) {
				return c.JSONPretty(
					http.StatusNotFound,
					EchoMessage{Msg: fmt.Sprintf("record with id=%d not found", id)},
					"\t",
				)
			}

			c.Logger().Errorf("failed to get observation: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, ob, "\t")
	}
}

func UpdateWeatherObservationHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ob models.WeatherObservation

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

		err = weatherService.UpdateWeatherObservation(c.Request().Context(), ob)
		if err != nil {
			if errors.As(err, &repository.ErrNotFound{}) {
				return c.JSONPretty(
					http.StatusNotFound,
					EchoMessage{Msg: fmt.Sprintf("record with id=%d not found", id)},
					"\t",
				)
			}

			c.Logger().Errorf("failed to update observation: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, EchoMessage{Msg: "successfully updated"}, "\t")
	}
}

func DeleteWeatherObservationHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseID(c)
		if err != nil {
			return c.JSONPretty(
				http.StatusBadRequest,
				EchoMessage{Msg: fmt.Sprintf("parseID: %s", err)},
				"\t",
			)
		}

		ob, err := weatherService.DeleteWeatherObservation(c.Request().Context(), id)
		if err != nil {
			if errors.As(err, &repository.ErrNotFound{}) {
				return c.JSONPretty(
					http.StatusNotFound,
					EchoMessage{Msg: fmt.Sprintf("record with id=%d not found", id)},
					"\t",
				)
			}

			c.Logger().Errorf("failed to delete observation: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		return c.JSONPretty(http.StatusOK, ob, "\t")
	}
}

func ListWeatherObservationsHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		obs, err := weatherService.ListWeatherObservations(c.Request().Context())
		if err != nil {
			c.Logger().Errorf("failed to get list of observations: %s", err)
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
