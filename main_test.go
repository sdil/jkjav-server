// Reference: https://github.com/gofiber/recipes/blob/master/unit-test/main_test.go

package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPPVRoute(t *testing.T) {
	tests := []struct {
		description string
		route string
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "list PWTC PPV",
			route:         "/list-ppv?state=PWTC",
			expectedCode:  200,
			expectedBody:  `[{"location":"PWTC","date":"20210516","availability":1000},{"location":"PWTC","date":"20210517","availability":1000},{"location":"PWTC","date":"20210518","availability":1000},{"location":"PWTC","date":"20210519","availability":1000},{"location":"PWTC","date":"20210520","availability":1000},{"location":"PWTC","date":"20210521","availability":1000},{"location":"PWTC","date":"20210522","availability":1000},{"location":"PWTC","date":"20210523","availability":1000},{"location":"PWTC","date":"20210524","availability":1000},{"location":"PWTC","date":"20210525","availability":1000},{"location":"PWTC","date":"20210526","availability":1000},{"location":"PWTC","date":"20210527","availability":1000},{"location":"PWTC","date":"20210528","availability":1000},{"location":"PWTC","date":"20210529","availability":1000},{"location":"PWTC","date":"20210530","availability":1000},{"location":"PWTC","date":"20210531","availability":1000},{"location":"PWTC","date":"20210601","availability":1000},{"location":"PWTC","date":"20210602","availability":1000},{"location":"PWTC","date":"20210603","availability":1000},{"location":"PWTC","date":"20210604","availability":1000},{"location":"PWTC","date":"20210605","availability":1000},{"location":"PWTC","date":"20210606","availability":1000},{"location":"PWTC","date":"20210607","availability":1000},{"location":"PWTC","date":"20210608","availability":1000},{"location":"PWTC","date":"20210609","availability":1000},{"location":"PWTC","date":"20210610","availability":1000},{"location":"PWTC","date":"20210611","availability":1000},{"location":"PWTC","date":"20210612","availability":1000},{"location":"PWTC","date":"20210613","availability":1000},{"location":"PWTC","date":"20210614","availability":1000}]`,
		},
		{
			description:   "list PWTC PPV without state param",
			route:         "/list-ppv",
			expectedCode:  400,
			expectedBody:  "Please select a state",
		},
	}

	// Setup the app as it is done in the main function
	app := Setup()

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		res, err := app.Test(req)

		// Ensure there's no error when requesting to the endpoint
		assert.Nilf(t, err, test.description)

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
