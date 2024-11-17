package models

import "time"

type WeatherObservation struct {
	ID            int       `json:"id"`
	Location      Location  `json:"location"`
	Timestamp     time.Time `json:"timestamp"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	Pressure      float64   `json:"pressure"`
	Wind          Wind      `json:"wind"`
	WeatherStatus string    `json:"weather_status"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

type Wind struct {
	Speed     float64 `json:"speed"`
	Direction float64 `json:"direction"`
}
