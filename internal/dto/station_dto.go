package dto

type StationRequest struct {
	StationID    string  `json:"station_id"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Location     string  `json:"location"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	BaseStationID   *uint   `json:"base_station_id"`
	InitialDistance float64 `json:"initial_distance"`
	URLStreaming    string  `json:"url_streaming"`
	SiteID          *uint   `json:"site_id"`
}

type StationResponse struct {
	ID              uint                 `json:"id"`
	StationID       string               `json:"station_id"`
	Name            string               `json:"name"`
	Latitude        float64              `json:"latitude"`
	Longitude       float64              `json:"longitude"`
	Location        string               `json:"location"`
	Status          string               `json:"status"`
	BaseStationID   *uint                `json:"base_station_id"`
	BaseStation     *BaseStationResponse `json:"base_station,omitempty"`
	InitialDistance float64              `json:"initial_distance"`
	URLStreaming    string               `json:"url_streaming"`
	SiteID          *uint                `json:"site_id"`
}
