package summary

type Data map[string]float64

type Summary struct {
	Id           string
	UserId       string
	ControllerId string
	Data         Data
}
