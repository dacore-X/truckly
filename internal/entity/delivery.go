package entity

import (
	"time"
)

// Delivery represents delivery data struct for internal use
type Delivery struct {
	ID        int       `json:"id"`
	ClientID  int       `json:"client_id"`
	CourierID int       `json:"courier_id"`
	StatusID  int       `json:"status_id"`
	TypeID    int       `json:"type_id"`
	Geo       *Geo      `json:"geo"`
	Price     float64   `json:"price"`
	HasLoader bool      `json:"has_loader"`
	CreatedAt time.Time `json:"created_at"`
}

// Geo represents geo data struct for internal use
type Geo struct {
	FromLongitude float64 `json:"from_longitude"`
	FromLatitude  float64 `json:"from_latitude"`
	FromObject    string  `json:"from_object"`
	ToLongitude   float64 `json:"to_longitude"`
	ToLatitude    float64 `json:"to_latitude"`
	ToObject      string  `json:"to_object"`
	Distance      float64 `json:"distance"`
}
