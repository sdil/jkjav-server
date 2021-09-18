package booking

import (
	"github.com/sdil/jkjav-server/pkg/entities"
)

type Service interface {
	InsertBooking(booking *entities.Booking) (entities.Booking, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) InsertBooking(booking *entities.Booking) (entities.Booking, error) {
	s.repository.CreateBooking(booking)

	return *booking, nil
}