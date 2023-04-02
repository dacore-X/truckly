package entity

import "time"

// Delivery represents the user's delivery data struct
type Delivery struct {
	ID            int       `json:"id"`
	ClientID      int       `json:"client_id"`
	CourierID     int       `json:"courier_id"`
	StatusID      int       `json:"status_id"`
	TruckID       int       `json:"truck_id"`
	FromLongitude float64   `json:"from_longitude"`
	FromLatitude  float64   `json:"from_latitude"`
	ToLongitude   float64   `json:"to_longitude"`
	ToLatitude    float64   `json:"to_latitude"`
	Distance      float64   `json:"distance"`
	Price         float64   `json:"price"`
	HasLoader     bool      `json:"has_loader"`
	CreatedAt     time.Time `json:"created_at"`
}
