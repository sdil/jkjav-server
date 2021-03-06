// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	entities "github.com/sdil/jkjav-server/pkg/entities"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CreateBooking provides a mock function with given fields: _a0
func (_m *Repository) CreateBooking(_a0 *entities.Booking) (*entities.Booking, error) {
	ret := _m.Called(_a0)

	var r0 *entities.Booking
	if rf, ok := ret.Get(0).(func(*entities.Booking) *entities.Booking); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Booking)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*entities.Booking) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublishMessage provides a mock function with given fields: _a0
func (_m *Repository) PublishMessage(_a0 *entities.Booking) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entities.Booking) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
