package dto

// DeliveryCreateRequestBody represents the request body with data
// sent by the user to API to create new delivery order
type DeliveryCreateRequestBody struct {
	TypeID    int           `json:"type_id" binding:"required,gte=1,lte=5"`
	FromPoint *PointRequest `json:"from_point" binding:"required"`
	ToPoint   *PointRequest `json:"to_point" binding:"required"`
	HasLoader bool          `json:"has_loader"`
}
