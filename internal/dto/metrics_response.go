package dto

// DeliveriesCntPerDay represents the response body
// with new and completed deliveries' counts per last 24 hours
type DeliveriesCntPerDay struct {
	NewCnt       int `json:"new_cnt"`
	CompletedCnt int `json:"completed_cnt"`
}

// RevenuePerDay represents the response body
// with revenue sum per last 24 hours
type RevenuePerDay struct {
	Revenue int `json:"revenue"`
}

// NewClientsCntPerDay represents the response body
// with new registered clients' count per last 24 hours
type NewClientsCntPerDay struct {
	NewClientsCnt int `json:"new_clients_cnt"`
}

// DeliveryTypesPercentPerDay represents the response body
// with different delivery types' percentages per last 24 hours
type DeliveryTypesPercentPerDay struct {
	FootPercent      float64 `json:"foot_percent"`
	CarPercent       float64 `json:"car_percent"`
	MinivanPercent   float64 `json:"minivan_percent"`
	TruckPercent     float64 `json:"truck_percent"`
	LongTruckPercent float64 `json:"long_truck_percent"`
}

// MetricsPerDayResponse represents the response body
// with all metrics per last 24 hours
type MetricsPerDayResponse struct {
	DeliveriesCnt        *DeliveriesCntPerDay        `json:"deliveries_cnt"`
	Revenue              int                         `json:"revenue"`
	NewClientsCnt        int                         `json:"new_clients_cnt"`
	DeliveryTypesPercent *DeliveryTypesPercentPerDay `json:"delivery_types_percent"`
}

// CurrentDeliveries represents the response body
// with brief information about current delivery
type CurrentDelivery struct {
	FromPoint *PointResponse `json:"from_point"`
	Price     float64        `json:"price"`
}

// MetricsDeliveriesResponse represents the response body
// with list of brief information about current deliveries
type MetricsDeliveriesResponse struct {
	Deliveries []*CurrentDelivery `json:"current_deliveries"`
}
