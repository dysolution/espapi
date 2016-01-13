package espsdk

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// A Result contains information relative to a completed request, including
// the time elapsed to fulfill the request and any errors.
type result struct {
	Response *response     `json:"response"`
	Payload  []byte        `json:"-"`
	Duration time.Duration `json:"response_time"`
	Err      error         `json:"-"`
}

func getResult(c *http.Client, req *http.Request) (*result, error) {
	httpCommand := req.Method + " " + string(req.URL.Path)
	start := start(httpCommand)
	resp, err := c.Do(req)
	duration := elapsed(httpCommand, start) / time.Millisecond
	if err != nil {
		log.Error(err)
		return buildResult(resp, nil, duration, err), err
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return buildResult(resp, payload, duration, err), nil
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.WithFields(log.Fields{
			"object":      "response",
			"status_code": resp.StatusCode,
			"status":      resp.Status,
		}).Warn()
	}
	return buildResult(resp, payload, duration, nil), nil
}

func buildResult(resp *http.Response, payload []byte, duration time.Duration, err error) *result {
	return &result{
		&response{
			resp.StatusCode,
			resp.Status,
		},
		payload, duration, err}
}
