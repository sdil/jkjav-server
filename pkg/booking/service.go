package booking

import (
	"github.com/sdil/jkjav-server/pkg/entities"
)

type Service interface {
	InsertBooking(booking *entities.Booking) (*entities.Booking, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) InsertBooking(booking *entities.Booking) (*entities.Booking, error) {
	
	// Insert the data in data store
	_, err := s.repository.CreateBooking(booking)
	if err != nil {
		return nil, err
	}

	// Publish message to Kafka broker
	err = s.repository.PublishMessage(booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}