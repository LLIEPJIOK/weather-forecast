// Code generated by mockery v2.45.1. DO NOT EDIT.

package repository_test

import (
	context "context"

	models "github.com/LLIEPJIOK/weather-forecast/backend/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

// AddWeather provides a mock function with given fields: ctx, weather
func (_m *MockDatabase) AddWeather(ctx context.Context, weather *models.Weather) (*models.Weather, error) {
	ret := _m.Called(ctx, weather)

	if len(ret) == 0 {
		panic("no return value specified for AddWeather")
	}

	var r0 *models.Weather
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Weather) (*models.Weather, error)); ok {
		return rf(ctx, weather)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Weather) *models.Weather); ok {
		r0 = rf(ctx, weather)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Weather)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Weather) error); ok {
		r1 = rf(ctx, weather)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteWeather provides a mock function with given fields: ctx, id
func (_m *MockDatabase) DeleteWeather(ctx context.Context, id int) (*models.Weather, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteWeather")
	}

	var r0 *models.Weather
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*models.Weather, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Weather); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Weather)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWeather provides a mock function with given fields: ctx, id
func (_m *MockDatabase) GetWeather(ctx context.Context, id int) (*models.Weather, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetWeather")
	}

	var r0 *models.Weather
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*models.Weather, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Weather); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Weather)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListWeathers provides a mock function with given fields: ctx
func (_m *MockDatabase) ListWeathers(ctx context.Context) ([]*models.Weather, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListWeathers")
	}

	var r0 []*models.Weather
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*models.Weather, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*models.Weather); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Weather)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateWeather provides a mock function with given fields: ctx, weather
func (_m *MockDatabase) UpdateWeather(ctx context.Context, weather *models.Weather) (*models.Weather, error) {
	ret := _m.Called(ctx, weather)

	if len(ret) == 0 {
		panic("no return value specified for UpdateWeather")
	}

	var r0 *models.Weather
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Weather) (*models.Weather, error)); ok {
		return rf(ctx, weather)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Weather) *models.Weather); ok {
		r0 = rf(ctx, weather)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Weather)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Weather) error); ok {
		r1 = rf(ctx, weather)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockDatabase creates a new instance of MockDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabase {
	mock := &MockDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
