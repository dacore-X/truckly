package dto

// DeliveryCreateBody represents the request body with data
// sent by the user to API to create new delivery order
type DeliveryCreateBody struct {
	TypeID    int           `json:"type_id" binding:"required,gte=1,lte=5"`
	FromPoint *PointRequest `json:"from_point" binding:"required"`
	ToPoint   *PointRequest `json:"to_point" binding:"required"`
	HasLoader bool          `json:"has_loader"`
}

// DeliveryIdURI represents URI with delivery's ID to get info
// of specific delivery
type DeliveryIdURI struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// DeliveryStatusChangeBody represents the request body for changing
// delivery status
type DeliveryStatusChangeBody struct {
	StatusID int `json:"status_id" binding:"required,gte=1,lte=4"`
}

type DeliveryListGeolocationQuery struct {
	Latitude  float64 `form:"lat" binding:"required"`
	Longitude float64 `form:"lon" binding:"required"`
	Page      int     `form:"page" binding:"required,min=1"`
}
