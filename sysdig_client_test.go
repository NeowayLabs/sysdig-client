// +build unit

package sysdig_client_test

import (
	"encoding/json"
	sysdigclient "github.com/NeowayLabs/sysdig-client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getSumMetricConfiguration() (string, []sysdigclient.Metric, sysdigclient.Period) {
	filter := "agent.tag.team = \"datapirates\" and bot.name = \"abi\""

	metrics := []sysdigclient.Metric{
		{
			Id: "bot.result.status.2xx.total",
			Aggregations: sysdigclient.Aggregation{
				Time:  "sum",
				Group: "sum",
			},
		},
	}
	period := sysdigclient.Period{}

	return filter, metrics, period
}

func getSumMetricResponseRecord() map[string]interface{} {
	return map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"d": []int{
					80171,
				},
			},
		},
		"start": 1568127600,
		"end":   1568214000,
	}
}

func getSuccessHttpHandler() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/data" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(getSumMetricResponseRecord())
		}
	}))
}

func TestGetSumMetricWithDaysPeriod(t *testing.T) {

	filter, metrics, period := getSumMetricConfiguration()
	period.Days = 30

	ts := getSuccessHttpHandler()
	defer ts.Close()

	s := sysdigclient.NewWithEndpoint(ts.URL)
	sumMetric, err := s.GetSumMetric(metrics, filter, period)

	assert.Nil(t, err)
	assert.Greater(t, sumMetric, 0, "Count should be greater than 0!")
}

func TestGetSumMetricWithHourPeriod(t *testing.T) {

	filter, metrics, period := getSumMetricConfiguration()
	period.Hours = 3

	ts := getSuccessHttpHandler()
	defer ts.Close()

	s := sysdigclient.NewWithEndpoint(ts.URL)
	sumMetric, err := s.GetSumMetric(metrics, filter, period)

	assert.Nil(t, err)
	assert.Greater(t, sumMetric, 0, "Count should be greater than 0!")
}

func TestGetSumMetricWithMinutesPeriod(t *testing.T) {

	filter, metrics, period := getSumMetricConfiguration()
	period.Minutes = 10

	ts := getSuccessHttpHandler()
	defer ts.Close()

	s := sysdigclient.NewWithEndpoint(ts.URL)
	sumMetric, err := s.GetSumMetric(metrics, filter, period)

	assert.Nil(t, err)
	assert.Greater(t, sumMetric, 0, "Count should be greater than 0!")
}

func TestGetSumMetricWithErrorInGetPeriodInSeconds(t *testing.T) {

	filter, metrics, period := getSumMetricConfiguration()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	s := sysdigclient.New()
	sumMetric, err := s.GetSumMetric(metrics, filter, period)

	assert.EqualError(t, err, "error on give a period in seconds: one of the period fields must be greater than zero", "Error should be equal!")
	assert.Equal(t, sumMetric, 0, "Sum metric must be zero!")
}

func TestGetSumMetricUnmarshalResponseError(t *testing.T) {

	filter, metrics, period := getSumMetricConfiguration()
	period.Minutes = 10

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/data" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]int{1, 2, 3})
		}
	}))
	defer ts.Close()

	s := sysdigclient.NewWithEndpoint(ts.URL)
	sumMetric, err := s.GetSumMetric(metrics, filter, period)

	assert.EqualError(t, err, "error on unmarshal response result: json: cannot unmarshal array into Go value of type sysdig_client.ApiResult", "Error should be equal!")
	assert.Equal(t, sumMetric, 0, "Sum metric must be zero!")
}

func TestGetSumMetricBlankDataResponseError(t *testing.T) {

	filter, metrics, period := getSumMetricConfiguration()
	period.Minutes = 10

	result := getSumMetricResponseRecord()
	delete(result, "data")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/data" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
		}
	}))
	defer ts.Close()

	s := sysdigclient.NewWithEndpoint(ts.URL)
	sumMetric, err := s.GetSumMetric(metrics, filter, period)

	assert.EqualError(t, err, "api returned the blank data field", "Error should be equal!")
	assert.Equal(t, sumMetric, 0, "Sum metric must be zero!")
}
