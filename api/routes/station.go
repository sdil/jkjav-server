package routes

import (
	"fmt"
	"time"

	"github.com/sdil/jkjav-server/pkg/station"

	"github.com/gofiber/fiber/v2"
)

func StationRouter(app fiber.Router, service station.Service) {
	app.Get("/stations/:name", listStations(service))
}

var (
	StartDate = time.Date(2021, time.May, 15, 0, 0, 0, 0, time.UTC)
	EndDate   = time.Date(2021, time.June, 15, 0, 0, 0, 0, time.UTC)
)

// ListStation godoc
// @Summary List Station
// @Description Get station slots by location
// @Accept  json
// @Produce  json
// @Param name path string true "select the location. The only available option is PWTC"
// @Success 200 {array} entities.Station
// @Router /stations/{name} [get]
func listStations(service station.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		location := c.Params("name")
		if location == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Please select a location")
		}

		stations, err := service.FetchStation("PWTC", StartDate, EndDate)
		if err != nil {
			c.SendString(err.Error())
		}

		// Set Cache-control header to 1s
		c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=1"))

		return c.JSON(stations)
	}
}
