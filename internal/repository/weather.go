package repository

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/LLIEPJIOK/weather-forecast/internal/models"
)

type WeatherRepository struct {
	mu     *sync.Mutex
	obs    map[int]models.WeatherObservation
	nextID int
}

func NewWeatherRepository() *WeatherRepository {
	return &WeatherRepository{
		mu:     &sync.Mutex{},
		obs:    make(map[int]models.WeatherObservation),
		nextID: 1,
	}
}

func (r *WeatherRepository) AddWeatherObservation(
	_ context.Context,
	ob models.WeatherObservation,
) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ob.ID = r.nextID
	r.obs[r.nextID] = ob

	r.nextID++

	return ob.ID, nil
}

func (r *WeatherRepository) GetWeatherObservation(
	_ context.Context,
	id int,
) (models.WeatherObservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ob, ok := r.obs[id]
	if !ok {
		return models.WeatherObservation{}, NewErrNotFound(id)
	}

	return ob, nil
}

func (r *WeatherRepository) UpdateWeatherObservation(
	_ context.Context,
	ob models.WeatherObservation,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.obs[ob.ID]
	if !ok {
		return NewErrNotFound(ob.ID)
	}

	r.obs[ob.ID] = ob
	return nil
}

func (r *WeatherRepository) DeleteWeatherObservation(
	_ context.Context,
	id int,
) (models.WeatherObservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ob, ok := r.obs[id]
	if !ok {
		return models.WeatherObservation{}, NewErrNotFound(id)
	}

	delete(r.obs, id)
	return ob, nil
}

func (r *WeatherRepository) ListWeatherObservations(
	_ context.Context,
) ([]models.WeatherObservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	obList := slices.SortedFunc(
		maps.Values(r.obs),
		func(first, second models.WeatherObservation) int {
			if first.ID < second.ID {
				return -1
			}

			if first.ID > second.ID {
				return 1
			}

			return 0
		},
	)

	return obList, nil
}
