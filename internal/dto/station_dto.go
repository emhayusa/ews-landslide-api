package dto

type StationRequest struct {
	StationID    string  `json:"station_id"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Location     string  `json:"location"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	HardwareType string  `json:"hardware_type"`
}

type StationResponse struct {
	ID           uint    `json:"id"`
	StationID    string  `json:"station_id"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Location     string  `json:"location"`
	Status       string  `json:"status"`
	HardwareType string  `json:"hardware_type"`
}
