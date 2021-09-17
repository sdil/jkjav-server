package station

import (
	"fmt"
	"time"

	"github.com/sdil/jkjav-server/pkg/entities"
)
type Service interface {
	InsertStation(station *entities.Station) (*entities.Station, error)
	FetchStations(location string, startDate time.Time, endDate time.Time) (*[]entities.Station, error)
}

type service struct {
	repository Repository
}
//NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) InsertStation(station *entities.Station) (*entities.Station, error) {
	return s.repository.CreateStation(station)
}

func (s *service) FetchStations(location string, startDate time.Time, endDate time.Time) (*[]entities.Station, error) {

	stations := []entities.Station{}

	// Iterate each date
	days := endDate.Sub(startDate).Hours() / 24
	daysInt := int(days)
	for i := 1; i < daysInt; i++ {
		date := startDate.Add(time.Hour * time.Duration(i) * time.Duration(24))
		dateString := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())

		station, err := s.repository.ReadStation(location, dateString)
		if err != nil {
			return &stations, err
		}

		stations = append(stations, *station)
	}
	return &stations, nil
}
