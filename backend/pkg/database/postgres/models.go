// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package postgres

import (
	"time"
)

type Weather struct {
	ID            int64
	Timestamp     time.Time
	City          string
	Country       string
	Temperature   float64
	Humidity      float64
	Pressure      float64
	WindSpeed     float64
	WeatherStatus string
}