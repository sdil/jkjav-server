package entities

type Booking struct {
	MySejahteraID string `json:"mysejahteraId" redis:"mysejahteraId" example:"900127015527"`
	FirstName     string `json:"firstName" redis:"firstName" example:"Fadhil"`
	LastName      string `json:"lastName" redis:"lastName" example:"Yaacob"`
	Address       string `json:"address" redis:"address" example:"Kuala Lumpur"`
	Location      string `json:"location" redis:"location" example:"PWTC"`
	PhoneNumber   string `json:"phoneNumber" redis:"phoneNumber" example:"0123456789"`
	Date          string `json:"date" redis:"date" example:"20210516"`
}