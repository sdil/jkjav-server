package routes

// import (
// 	service "github.com/sdil/jkjav-server/pkg/booking"
// 	"github.com/sdil/jkjav-server/pkg/entities"

// 	"github.com/gofiber/fiber/v2"
// )

// // Submit godoc
// // @Summary Create Booking Slot
// // @Description Create a vaccine booking slot
// // @Accept  json
// // @Produce  json
// // @Param booking body booking true "booking info"
// // @Success 200 {object} booking
// // @Router /booking [post]
// func CreateBooking(c *fiber.Ctx) error {
// 	booking := new(entities.Booking)
// 	if err := c.BodyParser(booking); err != nil {
// 		return err
// 	}

// 	err := service.InsertBooking(booking)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 	}

// 	// Publish message to Message Queue Broker

// 	return c.JSON(booking)
// }
