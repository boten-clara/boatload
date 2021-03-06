package apiTest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sMailund/boatload/src/core/applicationServices"
	"github.com/sMailund/boatload/src/core/domainEntities"
	"github.com/sMailund/boatload/src/external/http/api"
)

// TODO fix broken tests
// TODO make sure all json fields are copied to entity

type metServiceStub struct {
	submitMethod func() error
}

func (ms metServiceStub) SubmitData(_ domainEntities.TimeSeries) error {
	return ms.submitMethod()
}

func TestShouldMarshallPostBody(t *testing.T) {
	submitMethod := func() error {
		return nil
	}
	api.UploadService = createUploadStub(submitMethod)

	res, req := createResponseAndRequestStubs(testPayload)

	api.UploadTimeSeries(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expexted 200, got %v: %v\n", res.Code, res.Body)
	}
}

func TestShouldOnlyAcceptPostMethod(t *testing.T) {
	methods := []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPut,
		http.MethodTrace,
	}

	submitMethod := func() error {
		return nil
	}
	api.UploadService = createUploadStub(submitMethod)

	for _, method := range methods {
		res, req := createResponseAndRequestStubsWithMethod(testPayload, method)

		api.UploadTimeSeries(res, req)
		if res.Code != http.StatusMethodNotAllowed {
			t.Errorf("expexted 405 for method %v, got %v: %v\n", method, res.Code, res.Body)
		}
	}

}

func TestShouldOnlyRespond500OnMetServiceError(t *testing.T) {
	submitMethod := func() error {
		return errors.New("sample error")
	}
	api.UploadService = createUploadStub(submitMethod)

	res, req := createResponseAndRequestStubs(testPayload)

	api.UploadTimeSeries(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expexted 500, got %v: %v\n", res.Code, res.Body)
	}
}

func createUploadStub(submitMethod func() error) applicationServices.UploadService {
	serviceStub := struct{ metServiceStub }{}
	serviceStub.submitMethod = submitMethod

	return *applicationServices.CreateUploadService(serviceStub)
}

func createResponseAndRequestStubs(payload string) (*httptest.ResponseRecorder, *http.Request) {
	return createResponseAndRequestStubsWithMethod(payload, http.MethodPost)
}

func createResponseAndRequestStubsWithMethod(payload string, method string) (*httptest.ResponseRecorder, *http.Request) {
	body := []byte(payload)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(method, api.UploadRoute, bytes.NewReader(body))
	return res, req
}

const testPayload = `{
  "tstype": "test",
  "tseries": [
    {
      "header": {
        "id": {
          "gliderID": "testID",
          "paramter": "testparam"
        },
        "extra": {
          "source": "sadfasdf",
          "name": "stuff"
        }
      },
      "observations": [
        {
          "time": "2020-06-16T06:00:00Z",
          "body": {
            "pos": {
              "lon": "1",
              "lat": "2",
              "depth": "3",
              "qc_flag": "test"
            },
            "value": "123",
            "qc_flag": "test"
          }
        }
      ]
    }
  ]
}
`

func TestShouldMapCsvToTimeSeries(t *testing.T) {
	r := strings.NewReader(csvPayload)

	keys := []string{"temperature", "conductivity"}
	mapped, err := api.MapToSeriesObservation(keys, r)
	temperatureObservations := mapped["temperature"]
	conductivityObservations := mapped["conductivity"]

	if err != nil {
		t.Errorf("expected error to be nil, but got %v", err.Error())
	}

	if len(temperatureObservations) != 2 {
		t.Errorf("expected len(temperatureObservations) to be 2, but was %v", len(temperatureObservations))
	}

	if len(conductivityObservations) != 2 {
		t.Errorf("expected len(conductivityObservations) to be 2, but was %v", len(temperatureObservations))
	}

	if temperatureObservations[0].Time != "2021-09-18T14:40:40+02:00" {
		t.Errorf("expected first timestamp to be 2021-09-18T14:15:40+00:00, but was %v", temperatureObservations[0].Time)
	}

	if temperatureObservations[0].Body.Value != "69.69" {
		t.Errorf("expected first temperature to be 69.69, but was %v", temperatureObservations[0].Body.Value)
	}

	if conductivityObservations[1].Body.Value != "56.78" {
		t.Errorf("expected second conductivity to be 56.78, but was %v", temperatureObservations[0].Body.Value)
	}
}

const csvPayload = `timestamp,lat,lon,depth,temperature,conductivity
1631968840,58.144699,7.998280,33.33,69.69,420.69
1631969344,50.421478,8.593940,44.44,12.34,56.78`
