package dto

// GeoCoderResponse is a struct for JSON decoding response
type GeoCoderResponse struct {
	Result GeoCoderResult `json:"result"`
	Meta   Meta           `json:"meta"`
}

// GeoCoderResult is a struct for JSON decoding "Result" field
type GeoCoderResult struct {
	Items []Item `json:"items"`
}

// Meta is a struct for JSON decoding meta information of the http response
type Meta struct {
	StatusCode int `json:"code"`
}

// Item is a struct for JSON decoding "Item" field
type Item struct {
	Address  string        `json:"address_name"`
	FullName string        `json:"full_name"`
	Point    PointResponse `json:"point"`
}

// PointResponse is a struct of point Geo position with latitude and longitude
type PointResponse struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// DistanceResponse is a struct for JSON decoding response
type DistanceResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
	} `json:"routes"`
}
