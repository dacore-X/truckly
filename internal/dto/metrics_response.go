package dto

// DeliveriesCntPerDay represents the response body
// with new and completed deliveries' counts per last 24 hours
// and differences in percents between previous and current day
// for corresponding values
type DeliveriesCntPerDay struct {
	NewCnt           int     `json:"new_cnt"`
	NewCntDiff       float64 `json:"new_cnt_diff"`
	CompletedCnt     int     `json:"completed_cnt"`
	CompletedCntDiff float64 `json:"completed_cnt_diff"`
}

// RevenuePerDay represents the response body
// with revenue sum per last 24 hours and difference
// in percents between previous and current day for revenue
type RevenuePerDay struct {
	Revenue     int     `json:"revenue"`
	RevenueDiff float64 `json:"revenue_diff"`
}

// NewClientsCntPerDay represents the response body
// with new registered clients' count per last 24 hours and difference
// in percents between previous and current day for new registered clients' count
type NewClientsCntPerDay struct {
	NewClientsCnt     int     `json:"new_clients_cnt"`
	NewClientsCntDiff float64 `json:"new_clients_cnt_diff"`
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
	Revenue              *RevenuePerDay              `json:"revenue"`
	NewClientsCnt        *NewClientsCntPerDay        `json:"new_clients_cnt"`
	DeliveryTypesPercent *DeliveryTypesPercentPerDay `json:"delivery_types_percent"`
}

// CurrentDelivery represents the response body
// with brief information about current delivery
type CurrentDelivery struct {
	FromObject *GeoObjectResponse `json:"from_object"`
	ToObject   string             `json:"to_object"`
	Price      float64            `json:"price"`
}

// MetricsDeliveriesResponse represents the response body
// with list of brief information about current deliveries
type MetricsDeliveriesResponse struct {
	Deliveries []*CurrentDelivery `json:"current_deliveries"`
}
