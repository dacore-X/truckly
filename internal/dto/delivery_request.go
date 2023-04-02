package dto

// DeliveryCreateRequestBody represents the request body with data
// sent by the user to API to create new delivery order
type DeliveryCreateRequestBody struct {
	TruckID       int     `json:"truck_id"`
	FromLongitude float64 `json:"from_longitude" binding:"required"`
	FromLatitude  float64 `json:"from_latitude" binding:"required"`
	ToLongitude   float64 `json:"to_longitude" binding:"required"`
	ToLatitude    float64 `json:"to_latitude" binding:"required"`
	HasLoader     bool    `json:"has_loader"`
}
