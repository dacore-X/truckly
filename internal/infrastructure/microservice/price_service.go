package microservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/internal/dto"
	"io"
	"log"
	"net/http"
)

type PriceEstimator struct {
	ServicePort int
}

func New(cfg *config.SERVICES) *PriceEstimator {
	return &PriceEstimator{
		ServicePort: cfg.Ports["PriceEstimator"],
	}
}

// doRequest making request to URL in args and returns *http.Response
func doRequest(method, URL string, body io.Reader) (*http.Response, error) {
	//ctx := context.TODO()
	switch method {
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

func (p *PriceEstimator) EstimateDeliveryPrice(body *dto.EstimatePriceRequestBody) (float64, error) {
	URL := fmt.Sprintf("http://localhost:%v/price", p.ServicePort)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return 0, errors.New("error encoding body")
	}
	result, err := doRequest(http.MethodPost, URL, &buf)
	if result.StatusCode != 200 {
		return 0, errors.New("internal server error")
	}

	response := &dto.EstimatePriceResponse{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(response)

	result.Body.Close()

	if err != nil {
		//log.Println("error unmarshalling body")
		return 0, errors.New("error unmarshalling body")
	}

	return response.Price, nil
}
