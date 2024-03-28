package types

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type BusData struct {
	BusId string  `json:"busId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
	Index int     `json:"index"`
}

type UserData struct {
	UserId string  `json:"userId"`
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
}
