package models

import "time"

type WeatherObservation struct {
	ID            int       `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	City          string    `json:"city"`
	Country       string    `json:"country"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	Pressure      float64   `json:"pressure"`
	WindSpeed     float64   `json:"wind_speed"`
	WeatherStatus string    `json:"weather_status"`
}
