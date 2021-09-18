package booking

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/sdil/jkjav-server/pkg/entities"
	"gopkg.in/Shopify/sarama.v1"
)

type Repository interface {
	CreateBooking(booking *entities.Booking) (*entities.Booking, error)
	PublishMessage(booking *entities.Booking) error
}

type repository struct {
	Pool          *redis.Pool
	MessageBroker sarama.SyncProducer
}

func NewRepo(pool *redis.Pool, mb sarama.SyncProducer) Repository {
	return &repository{
		Pool:          pool,
		MessageBroker: mb,
	}
}

func (r *repository) CreateBooking(booking *entities.Booking) (*entities.Booking, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	stationKey := "location:" + booking.Location + ":" + booking.Date

	ok, err := redis.Bool(conn.Do("EXISTS", stationKey))
	if err != nil {
		return nil, fmt.Errorf("failed to location key %s. error: %v", stationKey, err)
	}
	if ok == false {
		return nil, fmt.Errorf("station %s location & date combination is invalid", stationKey)
	}

	log.Printf("adding new booking %s and updating station %s availability", booking.MySejahteraID, stationKey)

	// Start a transaction
	if _, err = conn.Do("WATCH", stationKey); err != nil {
		return nil, fmt.Errorf("Failed to watch key %s: %v", stationKey, err)
	}

	var availability []byte
	availability, err = redis.Bytes(conn.Do("GET", stationKey))
	if err != nil {
		return nil, fmt.Errorf("error getting key %s: %v", stationKey, err)
	}

	log.Printf("%s availability is %s. attempting to add booking %s", stationKey, string(availability), booking.MySejahteraID)

	if string(availability) == "0" {
		if _, err = conn.Do("UNWATCH"); err != nil {
			log.Printf("%s failed to unwatch", stationKey)
			return nil, fmt.Errorf("Failed to unwatch key %s: %v", stationKey, err)
		}

		log.Printf("%s is fully booked", stationKey)
		return nil, fmt.Errorf("Sorry, station is fully booked")
	}

	conn.Send("MULTI")

	// Decrease the counter
	if err := conn.Send("DECR", stationKey); err != nil {
		return nil, fmt.Errorf("error setting key %s: %v", stationKey, err)
	}

	// Add new booking dict
	bookingKey := "booking:" + booking.MySejahteraID
	if err := conn.Send("HMSET", redis.Args{}.Add(bookingKey).AddFlat(booking)...); err != nil {
		return nil, fmt.Errorf("error setting key %s: %v", bookingKey, err)
	}

	// Execute Transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		return nil, fmt.Errorf("error setting key %s: %v", bookingKey, err)
	}

	log.Printf("successfully added new booking %s and updated station %s availability", booking.MySejahteraID, stationKey)

	return booking, nil
}

func (r *repository) PublishMessage(booking *entities.Booking) error {

	if r.MessageBroker != nil {
		r.MessageBroker.SendMessage(&sarama.ProducerMessage{
			Topic: "booking-slot",
			Value: sarama.StringEncoder("test"),
		})
		return nil
	} else {
		log.Println("Silently ignore this error. Your infra probably don't have Kafka available.")
		return nil
	}
}
