package booking

import (
	"github.com/sdil/jkjav-server/pkg/entities"
)

func InsertBooking(booking *entities.Booking) error {
	conn := Pool.Get()
	defer conn.Close()

	ppvKey := "location:" + booking.Location + ":" + booking.Date

	ok, err := redis.Bool(conn.Do("EXISTS", ppvKey))
	if err != nil {
		return fmt.Errorf("failed to location key %s. error: %v", ppvKey, err)
	}
	if ok == false {
		return fmt.Errorf("ppv %s location & date combination is invalid", ppvKey)
	}

	log.Printf("adding new booking %s and updating ppv %s availability", booking.MySejahteraID, ppvKey)

	// Start a transaction
	if _, err = conn.Do("WATCH", ppvKey); err != nil {
		return fmt.Errorf("Failed to watch key %s: %v", ppvKey, err)
	}

	var availability []byte
	availability, err = redis.Bytes(conn.Do("GET", ppvKey))
	if err != nil {
		return fmt.Errorf("error getting key %s: %v", ppvKey, err)
	}

	log.Printf("%s availability is %s. attempting to add booking %s", ppvKey, string(availability), booking.MySejahteraID)

	if string(availability) == "0" {
		if _, err = conn.Do("UNWATCH"); err != nil {
			log.Printf("%s failed to unwatch", ppvKey)
			return fmt.Errorf("Failed to unwatch key %s: %v", ppvKey, err)
		}

		log.Printf("%s is fully booked", ppvKey)
		return fmt.Errorf("Sorry, ppv is fully booked")
	}

	conn.Send("MULTI")

	// Decrease the counter
	if err := conn.Send("DECR", ppvKey); err != nil {
		return fmt.Errorf("error setting key %s: %v", ppvKey, err)
	}

	// Add new booking dict
	bookingKey := "booking:" + booking.MySejahteraID
	if err := conn.Send("HMSET", redis.Args{}.Add(bookingKey).AddFlat(booking)...); err != nil {
		return fmt.Errorf("error setting key %s: %v", bookingKey, err)
	}

	// Execute Transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		return fmt.Errorf("error setting key %s: %v", bookingKey, err)
	}

	log.Printf("successfully added new booking %s and updated ppv %s availability", booking.MySejahteraID, ppvKey)

	return nil
}