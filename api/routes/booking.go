package routes

import (
	"github.com/sdil/jkjav-server/pkg/booking"
	"github.com/sdil/jkjav-server/pkg/entities"

	"github.com/gofiber/fiber/v2"
)

func BookingRouter(app fiber.Router, service booking.Service) {
	app.Post("/booking", createBooking(service))
}

// Submit godoc
// @Summary Create Booking Slot
// @Description Create a vaccine booking slot
// @Accept  json
// @Produce  json
// @Param booking body entities.Booking true "booking info"
// @Success 200 {object} entities.Booking
// @Router /booking [post]
func createBooking(service booking.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		booking := new(entities.Booking)
		if err := c.BodyParser(booking); err != nil {
			return err
		}
	
		_, err := service.InsertBooking(booking)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	
		return c.JSON(booking)
	}
}
