package dto

// DistanceRequest represents the request body with data
// sent by the server to 2GIS API to find distance between points
type DistanceRequest struct {
	Points  []PointRequest `json:"points"`
	Sources []int          `json:"sources"`
	Targets []int          `json:"targets"`
	Type    string         `json:"type"`
}

type PointRequest struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
