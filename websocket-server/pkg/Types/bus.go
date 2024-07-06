package types

type BusData struct {
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
	Index int     `json:"index"`
}

type BusClose struct {
	Close bool `json:"close"`
}

type BusRequest struct {
	BusId    string   `json:"busId"`
	BusData  BusData  `json:"busData"`
	BusClose BusClose `json:"busClose"`
	Which    string   `json:"which"`
}

type BusMessage struct {
	Message string `json:"message"`
}

type BusUser struct {
	UserId string `json:"userId"`
}

type BusUserLocation struct {
	UserId string  `json:"userId"`
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
}

type BusResponse struct {
	BusId           string          `json:"busId"`
	BusData         BusData         `json:"busData"`
	BusMessage      BusMessage      `json:"busMessage"`
	BusUserList     []string        `json:"busUserList"`
	BusUserLocation BusUserLocation `json:"BusUserLocation"`
	Which           string          `json:"which"`
}
