package types

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type BusResponse struct {
	BusData  BusData  `json:"busData"`
	BusClose BusClose `json:"busClose"`
	Which    string   `json:"which"`
}

type BusData struct {
	BusId string  `json:"busId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
	Index int     `json:"index"`
}

type BusClose struct {
	Close bool `json:"close"`
}

type UserData struct {
	UserId string  `json:"userId"`
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
}
