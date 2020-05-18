package summary

type Summary struct {
	Id           string
	UserId       string
	ControllerId string
	Date         string
	Temperature  float64
	Humidity     float64
	Light        float64
	SoilMoisture int
	WaterLevel   int
}
