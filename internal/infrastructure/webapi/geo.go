package webapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/dacore-x/truckly/config"

	"github.com/dacore-x/truckly/internal/dto"
)

// Geo is a struct for communicating with 2GIS API
type Geo struct {
	BaseURLCatalog string
	BaseURLRouting string
	APIKeys        map[string]string
}

// URLQuery is a struct for building request URL
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
		r.Header.Set("Content-Type", "application/json")
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

// buildQuery building URL for request with input URLQuery
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

// GetCoordsByObject converts query to object dto.PointResponse
func (g *Geo) GetCoordsByObject(q string) (*dto.PointResponse, error) {
	if q == "" {
		return nil, errors.New("query is empty")
	}

	u := &URLQuery{
		base:     g.BaseURLCatalog,
		endpoint: "/3.0/items/geocode",
		params: map[string]string{
			"q":       q,
			"key":     g.APIKeys["catalog"],
			"fields":  "items.point",
			"city_id": "4504222397630173",
		},
	}

	URL := buildQuery(u)
	result, err := doRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	response := &dto.GeoCoderResponse{}
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
		log.Println("results not found by query")
		return nil, errors.New("results not found by query")
	}

	// returning only the first result
	return &response.Result.Items[0].Point, nil
}

// GetObjectByCoords converts geo object with input latitude and longitude to string representation
func (g *Geo) GetObjectByCoords(lat, lon float64) (string, error) {
	if lat == 0 || lon == 0 {
		return "", errors.New("coordinate couldn't be zero")
	}

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
		return "", err
	}

	response := &dto.GeoCoderResponse{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(response)
	result.Body.Close()

	if err != nil {
		log.Println("error unmarshalling meta")
		return "", errors.New("error unmarshalling meta")
	}

	if response.Meta.StatusCode >= 400 {
		log.Println("bad status code from geo")
		return "", errors.New("bad status code from geo")
	}

	if len(response.Result.Items) == 0 {
		log.Println("results not found by query")
		return "", errors.New("results not found by query")
	}
	// returning only the first result
	addr := response.Result.Items[0].Address
	if addr == "" {
		return response.Result.Items[0].FullName, nil
	}
	return response.Result.Items[0].Address, nil
}

// GetDistanceBetweenPoints calculating distance between 2 points (from and to) with input latitude and longitude
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
	body := dto.DistanceRequest{
		Points: []dto.PointRequest{
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
	result, err := doRequest(http.MethodPost, URL, &buf)
	if result.StatusCode != 200 {
		return 0, errors.New("error response 2gis")
	}

	response := &dto.DistanceResponse{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(response)
	result.Body.Close()

	if err != nil {
		//log.Println("error unmarshalling body")
		return 0, errors.New("error unmarshalling body")
	}

	if len(response.Routes) == 0 {
		return 0, errors.New("routes not found")
	}
	return response.Routes[0].Distance, nil
}
