package types

type UserRequest struct {
	UserLocation Location `json:"userLocation"`
	Which        string   `json:"which"`
}

type UserBusData struct {
	BusId string  `json:"busId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type UserMessage struct {
	Message string `json:"message"`
}

// type UserData struct {
// 	UserId string  `json:"userId"`
// 	Lat    float64 `json:"lat"`
// 	Long   float64 `json:"long"`
// }

type UserResponse struct {
	UserId      string      `json:"userId"`
	UserBusData UserBusData `json:"userBusData"`
	UserMessage UserMessage `json:"userMessage"`
	Which       string      `json:"which"`
}
