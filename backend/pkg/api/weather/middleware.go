package weather

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, err := c.Cookie("X-User-Role")
		if err != nil {
			c.Logger().Errorf("failed to get cookie: %s", err)
			return c.JSONPretty(
				http.StatusInternalServerError,
				EchoMessage{Msg: "The server is temporarily unavailable, please try again later"},
				"\t",
			)
		}

		fmt.Println(role.Value)

		if role.Value != "admin" {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}
