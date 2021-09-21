package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	tests := []struct {
		description  string
		route        string
		booking      *entities.Booking
		expectedCode int
		expectedBody string
	}{
		{
			description:  "Create valid booking",
			route:        "/booking",
			booking:      &entities.Booking{MySejahteraID: "900127015527", FirstName: "Fadhil", LastName: "Yaacob", Location: "PWTC", Address: "Kuala Lumpur", PhoneNumber: "0123456789", Date: "20210517"},
			expectedCode: 200,
			expectedBody: `{"mysejahteraId":"900127015527","firstName":"Fadhil","lastName":"Yaacob","address":"Kuala Lumpur","location":"PWTC","phoneNumber":"0123456789","date":"20210517"}`,
		},
		{
			description:  "Create valid booking",
			route:        "/booking",
			booking:      &entities.Booking{MySejahteraID: "", FirstName: "Fadhil", LastName: "Yaacob", Location: "PWTC", Address: "Kuala Lumpur", PhoneNumber: "0123456789", Date: "20210517"},
			expectedCode: 400,
			expectedBody: `{"message":"MySejahteraID cannot be empty","status":"Failed"}`,
		},
	}

	app := fiber.New()
	BookingRouter(app, bookingService)

	for _, test := range tests {

		bookingRepository.On("CreateBooking", test.booking).Return(test.booking, nil)
		bookingRepository.On("PublishMessage", test.booking).Return(nil)

		// Create a new http request with the route from the test case
		jsonBooking, _ := json.Marshal(test.booking)
		req := httptest.NewRequest("POST", test.route, bytes.NewBuffer(jsonBooking))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, 1)

		assert.Nilf(t, err, fmt.Sprintf("Failed to make request to %s", test.route))

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, fmt.Sprintf("'%s' test status code is not same", test.description))

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		assert.Nilf(t, err, fmt.Sprintf("Failed to read response body %s", test.route))

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, fmt.Sprintf("The response body %s is empty", test.route))

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), fmt.Sprintf("'%s' test body is not same", test.description))
	}
}
