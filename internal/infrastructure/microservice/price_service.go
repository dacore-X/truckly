package microservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
)

type PriceEstimator struct {
	ServicePort int
	appLogger   *logger.Logger
}

func New(cfg *config.SERVICES, l *logger.Logger) *PriceEstimator {
	return &PriceEstimator{
		ServicePort: cfg.Ports["PriceEstimator"],
		appLogger:   l,
	}
}

// doRequest making request to URL in args and returns *http.Response
func doRequest(method, URL string, body io.Reader) (*http.Response, error) {
	switch method {
	case http.MethodPost:
		r, err := http.Post(URL, "application/json", body)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, errors.New("error making request")
}

// EstimateDeliveryPrice making request to PriceEstimator service to get price for delivery
func (p *PriceEstimator) EstimateDeliveryPrice(body *dto.EstimatePriceInternalRequestBody) (float64, error) {
	URL := fmt.Sprintf("http://localhost:%v/price", p.ServicePort)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		err := errors.New("error encoding body")
		p.appLogger.Error(err)
		return 0, err
	}
	result, err := doRequest(http.MethodPost, URL, &buf)
	if err != nil {
		p.appLogger.Errorf("microservice.doRequest: %v", err)
		return 0, err
	}

	if result.StatusCode != 200 {
		err := errors.New("internal server error")
		p.appLogger.Error(err)
		return 0, err
	}

	response := &dto.EstimatePriceResponse{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(response)

	result.Body.Close()

	if err != nil {
		err := errors.New("error unmarshalling body")
		p.appLogger.Error(err)
		return 0, err
	}

	return response.Price, nil
}
