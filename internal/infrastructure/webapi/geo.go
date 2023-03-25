package webapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dacore-x/truckly/config"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Geo is a struct for communicating with 2GIS API
type Geo struct {
	BaseURLCatalog string
	BaseURLRouting string
	APIKeys        map[string]string
}

// Meta is a struct for JSON decoding meta information of the http response
type Meta struct {
	StatusCode int `json:"code"`
}

// GeoCoderResponse is a struct for JSON decoding response
type GeoCoderResponse struct {
	Result GeoCoderResult `json:"result"`
	Meta   Meta           `json:"meta"`
}

// DistanceResponse is a struct for JSON decoding response
type DistanceResponse struct {
	Routes []Route `json:"routes"`
}

// Route is a struct for JSON decoding response
type Route struct {
	Distance float64 `json:"distance"`
}

// GeoCoderResult is a struct for JSON decoding "Result" field
type GeoCoderResult struct {
	Items []Item `json:"items"`
}

// Item is a struct for JSON decoding "Item" field
type Item struct {
	Address Address `json:"address_name"`
	Point   Point   `json:"point"`
}

// Address is a struct for JSON decoding "Address" field
type Address string

// Point is a struct for JSON decoding "Point" field
type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type DistanceRequest struct {
	Points  []Point `json:"points"`
	Sources []int   `json:"sources"`
	Targets []int   `json:"targets"`
	Type    string  `json:"type"`
}

type URLQuery struct {
	base     string
	endpoint string
	params   map[string]string
}

func New(cfg *config.GEO) *Geo {
	return &Geo{
		BaseURLCatalog: cfg.BaseURLCatalog,
		BaseURLRouting: cfg.BaseURLRouting,
		APIKeys: map[string]string{
			"catalog":    cfg.APIKeyCatalog,
			"navigation": cfg.APIKeyRouting,
		},
	}
}

// doRequest making request to URL in args and returns *http.Response
func doRequest(method, URL string, body io.Reader) (*http.Response, error) {
	//ctx := context.TODO()
	switch method {
	case http.MethodGet:
		r, err := http.Get(URL)
		//r.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Println("error creating request")
			return nil, err
		}
		return r, nil

	case http.MethodPost:
		r, err := http.Post(URL, "application/json", body)
		if err != nil {
			log.Println("error creating request")
			return nil, err
		}
		return r, nil
	}
	return nil, errors.New("error making request")
}

func buildQuery(u *URLQuery) string {
	URL, _ := url.Parse(u.base)
	URL.Path += u.endpoint
	params := url.Values{}
	for key, value := range u.params {
		params.Add(key, value)
	}
	URL.RawQuery = params.Encode()
	return URL.String()
}

func (g *Geo) GetCoordsByObject(q string) (*Point, error) {
	if q == "" {
		return nil, errors.New("query is empty")
	}

	u := &URLQuery{
		base:     g.BaseURLCatalog,
		endpoint: "/3.0/items/geocode",
		params: map[string]string{
			"q":      q,
			"key":    g.APIKeys["catalog"],
			"fields": "items.point",
		},
	}

	URL := buildQuery(u)
	result, err := doRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	response := &GeoCoderResponse{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(response)
	result.Body.Close()

	if err != nil {
		log.Println("error unmarshalling meta")
		return nil, errors.New("error unmarshalling meta")
	}

	if response.Meta.StatusCode >= 400 {
		log.Println("bad status code from geo")
		return nil, errors.New("bad status code from geo")
	}

	if len(response.Result.Items) == 0 {
		return nil, errors.New("results not found by query")
	}

	// returning only the first result
	return &response.Result.Items[0].Point, nil
}

func (g *Geo) GetObjectByCoords(lat, lon float64) (*Address, error) {
	u := &URLQuery{
		base:     g.BaseURLCatalog,
		endpoint: "/3.0/items/geocode",
		params: map[string]string{
			"lon":    fmt.Sprint(lon),
			"lat":    fmt.Sprint(lat),
			"key":    g.APIKeys["catalog"],
			"fields": "items.point",
		},
	}

	URL := buildQuery(u)
	result, err := doRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	response := &GeoCoderResponse{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(response)
	result.Body.Close()

	if err != nil {
		log.Println("error unmarshalling meta")
		return nil, errors.New("error unmarshalling meta")
	}

	if response.Meta.StatusCode >= 400 {
		log.Println("bad status code from geo")
		return nil, errors.New("bad status code from geo")
	}

	if len(response.Result.Items) == 0 {
		return nil, errors.New("results not found by query")
	}
	// returning only the first result
	return &response.Result.Items[0].Address, nil
}

func (g *Geo) GetDistanceBetweenPoints(latFrom, lonFrom, latTo, lonTo float64) (float64, error) {
	if latFrom == 0 || lonFrom == 0 || latTo == 0 || lonTo == 0 {
		return 0, errors.New("coordinate couldn't be zero")
	}

	u := &URLQuery{
		base:     g.BaseURLRouting,
		endpoint: "/get_dist_matrix",
		params: map[string]string{
			"key":     g.APIKeys["navigation"],
			"version": "2.0",
		},
	}
	URL := buildQuery(u)
	log.Println(URL)
	body := DistanceRequest{
		Points: []Point{
			{Lat: latFrom, Lon: lonFrom},
			{Lat: latTo, Lon: lonTo},
		},
		Sources: []int{0},
		Targets: []int{1},
		Type:    "jam",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return 0, errors.New("error encoding body")
	}
	log.Println(buf.String())
	result, err := doRequest(http.MethodPost, URL, &buf)
	log.Println(result)
	if result.StatusCode != 200 {
		return 0, errors.New("error response 2gis")
	}

	response := &DistanceResponse{}
	log.Println(result.Body)
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(response)
	result.Body.Close()

	if err != nil {
		log.Println("error unmarshalling body")
		return 0, errors.New("error unmarshalling body")
	}

	if len(response.Routes) == 0 {
		return 0, errors.New("routes not found")
	}

	return response.Routes[0].Distance, nil
}
