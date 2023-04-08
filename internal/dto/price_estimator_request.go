package dto

import "time"

// EstimatePriceInternalRequestBody represents the request body for internal
// request to Delivery Price Estimator Service
type EstimatePriceInternalRequestBody struct {
	TypeID    int       `json:"type_id" binding:"required"`
	HasLoader bool      `json:"has_loader"`
	Time      time.Time `json:"time" binding:"required"`
	Distance  float64   `json:"distance" binding:"required"` // km
}

// EstimatePriceRequestBody represents the request body with data
// sent by the user to API to estimate delivery price
type EstimatePriceRequestBody struct {
	TypeID    int           `json:"type_id" binding:"required,gte=1,lte=5"`
	FromPoint *PointRequest `json:"from_point" binding:"required"`
	ToPoint   *PointRequest `json:"to_point" binding:"required"`
	HasLoader bool          `json:"has_loader"`
}
