package dto

type BaseStationRequest struct {
	Kode   string  `json:"kode" validate:"required"`
	Nama   string  `json:"nama" validate:"required"`
	Lokasi string  `json:"lokasi"`
	Long   float64 `json:"long"`
	Lat    float64 `json:"lat"`
	SiteID *uint   `json:"site_id"`
}

type BaseStationResponse struct {
	ID        uint    `json:"id"`
	UUID      string  `json:"uuid"`
	Kode      string  `json:"kode"`
	Nama      string  `json:"nama"`
	Lokasi    string  `json:"lokasi"`
	Long      float64 `json:"long"`
	Lat       float64 `json:"lat"`
	SiteID    *uint   `json:"site_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
