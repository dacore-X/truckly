package dto

import "time"

type EstimatePriceRequestBody struct {
	TypeID    int       `json:"type_id" binding:"required"`
	HasLoader bool      `json:"has_loader"`
	Time      time.Time `json:"time" binding:"required"`
	Distance  float64   `json:"distance" binding:"required"` // km
}
