package entities

type Station struct {
	Location     string `json:"location" redis:"location" example:"PWTC"`
	Date         string `json:"date" redis:"date" example:"20210516"`
	Availability int    `json:"availability" redis:"availability" example:"10"`
}
