package dto

type SiteRequest struct {
	Nama       string `json:"nama" validate:"required"`
	Lokasi     string `json:"lokasi"`
	Keterangan string `json:"keterangan"`
}

type SiteResponse struct {
	ID           uint                  `json:"id"`
	Nama         string                `json:"nama"`
	Lokasi       string                `json:"lokasi"`
	Keterangan   string                `json:"keterangan"`
	BaseStations []BaseStationResponse `json:"base_stations,omitempty"`
	Stations     []StationResponse     `json:"stations,omitempty"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}
