package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/db"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/models"
)

type Location struct {
	DisplayName string `json:"display_name"`
}

func UpdateUserLocationDB(userId string, lat float64, long float64) {
	apiUrl := "https://nominatim.openstreetmap.org/reverse.php?lat=" + strconv.FormatFloat(lat, 'f', 2, 64) + "&lon=" + strconv.FormatFloat(long, 'f', 2, 64) + "&zoom=18&format=jsonv2"
	client := &http.Client{}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var location Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		fmt.Println("Error in Parsing the Body")
		return
	}

	if result := db.DB.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"lat":      lat,
		"long":     long,
		"location": location.DisplayName,
	}); result.Error != nil {
		fmt.Println("Update User Failed")
	}
}
