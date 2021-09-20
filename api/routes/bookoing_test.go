package routes

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/sdil/jkjav-server/pkg/booking"
	"github.com/sdil/jkjav-server/pkg/booking/mocks"
	"github.com/sdil/jkjav-server/pkg/entities"
	"github.com/stretchr/testify/assert"
)

var bookingRepository = &mocks.Repository{}
var bookingService = booking.NewService(bookingRepository)

func TestCreate(t *testing.T) {
	booking := &entities.Booking{MySejahteraID: "900127015527", FirstName: "Fadhil", LastName: "Yaacob", Location: "PWTC", Address: "Kuala Lumpur", PhoneNumber: "0123456789", Date: "20210516"}
	bookingRepository.On("CreateBooking", booking).Return(booking, nil)
	bookingRepository.On("PublishMessage", booking).Return(nil)

	app := fiber.New()
	BookingRouter(app, bookingService)

	// Create a new http request with the route from the test case
	var jsonStr = []byte(`{"address":"Kuala Lumpur","date":"20210516","firstName":"Fadhil","lastName":"Yaacob","location":"PWTC","mysejahteraId":"900127015527","phoneNumber":"0123456789"}`)
	req := httptest.NewRequest("POST", "/booking", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "OK")
}
