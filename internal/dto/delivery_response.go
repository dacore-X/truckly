package dto

import "time"

type DeliveryFullInfoResponse struct {
	ID         int                 `json:"id"`
	TypeID     int                 `json:"type_id"`
	Courier    DeliveryCourierInfo `json:"courier"`
	StatusID   int                 `json:"status_id"`
	Price      float64             `json:"price"`
	HasLoader  bool                `json:"has_loader"`
	FromObject GeoObjectResponse   `json:"from_object"`
	ToObject   GeoObjectResponse   `json:"to_object"`
	Distance   int                 `json:"distance"`
	Time       time.Time           `json:"time"`
}

type DeliveryBriefResponse struct {
	ID         int       `json:"id"`
	TypeID     int       `json:"type_id"`
	HasLoader  bool      `json:"has_loader"`
	StatusID   int       `json:"status_id"`
	Price      float64   `json:"price"`
	FromObject string    `json:"from_object"`
	ToObject   string    `json:"to_object"`
	Distance   int       `json:"distance"`
	Time       time.Time `json:"time"`
}

type GeoObjectResponse struct {
	Object    string  `json:"object"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DeliveryCourierInfo struct {
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}
