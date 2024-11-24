CREATE TABLE IF NOT EXISTS weather
(
  id             BIGINT           NOT NULL GENERATED ALWAYS AS IDENTITY,
  timestamp      timestamp        NOT NULL,
  city           TEXT             NOT NULL,
  country        TEXT             NOT NULL,
  temperature    double precision NOT NULL,
  humidity       double precision NOT NULL,
  pressure       double precision NOT NULL,
  wind_speed     double precision NOT NULL,
  weather_status TEXT             NOT NULL,
  PRIMARY KEY (id)
);