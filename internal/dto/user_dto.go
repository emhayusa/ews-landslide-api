package dto

import "time"

type UserRequest struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password,omitempty"`
	FullName string   `json:"full_name"`
	Role     string   `json:"role"`
	SiteIDs  []uint   `json:"site_ids"`
}

type UserResponse struct {
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	FullName  string         `json:"full_name"`
	Role      string         `json:"role"`
	Sites     []SiteResponse `json:"sites,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}
