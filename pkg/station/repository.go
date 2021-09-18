package station

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/sdil/jkjav-server/pkg/entities"
)

//Repository interface allows us to access the CRUD Operations in mongo here.
type Repository interface {
	CreateStation(station *entities.Station) (*entities.Station, error)
	ReadStation(location string, date string) (*entities.Station, error)
}
type repository struct {
	Pool *redis.Pool
}
//NewRepo is the single instance repo that is being created.
func NewRepo(pool *redis.Pool) Repository {
	return &repository{
		Pool: pool,
	}
}

func (r *repository) CreateStation(station *entities.Station) (*entities.Station, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	key := "location:" + station.Location + ":" + station.Date

	_, err := conn.Do("SET", key, station.Availability)
	if err != nil {
		return station, fmt.Errorf("error setting key %s: %v", key, err)
	}

	return station, err
}

func (r *repository) ReadStation(location string, date string) (*entities.Station, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	key := "location:" + location + ":" + date

	var station entities.Station
	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, fmt.Errorf("error getting key %s: %v", key, err)
	}

	availability, err := strconv.Atoi(string(data))
	if err != nil {
		log.Println("error reading location availability: " + err.Error())
	}

	station.Location = location
	station.Date = date
	station.Availability = availability

	return &station, err
}