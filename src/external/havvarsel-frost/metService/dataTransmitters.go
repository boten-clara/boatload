package metService

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type iDataTransmitter interface {
	SendData(r *http.Request) error
}

type devDataTransmitter struct {
}

func (d devDataTransmitter) SendData(r *http.Request) error {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	fmt.Printf("intercepted message:\n%v\n\n", string(requestDump))
	return nil
}

type prodDataTransmitter struct {
}

// TODO: this implementation is untested
func (p prodDataTransmitter) SendData(r *http.Request) error {
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("request got non success status code: %v, %v", res.StatusCode, res.Status))
	}

	return nil
}
