-- name: AddWeather :one
INSERT INTO weather (timestamp, temperature, humidity, pressure, wind_speed, city, country, weather_status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetWeather :one
SELECT * 
FROM weather
WHERE id = $1;

-- name: UpdateWeather :one
UPDATE weather
SET 
    timestamp = $2,
    temperature = COALESCE(NULLIF($3, 0), temperature),
    humidity = COALESCE(NULLIF($4, 0), humidity),
    pressure = COALESCE(NULLIF($5, 0), pressure),
    wind_speed = COALESCE(NULLIF($6, 0), wind_speed),
    city = COALESCE(NULLIF($7, ''), city),
    country = COALESCE(NULLIF($8, ''), country),
    weather_status = COALESCE(NULLIF($9, ''), weather_status)
WHERE id = $1
RETURNING *;

-- name: DeleteWeather :one
DELETE FROM weather
WHERE id = $1
RETURNING *;

-- name: ListWeathers :many
SELECT * 
FROM weather;
