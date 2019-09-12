package sysdig_client

import (
	"encoding/json"
	"fmt"
	"github.com/NeowayLabs/sysdig-client/client"
)

type Sysdigclient struct {
	client client.Client
}

type Query struct {
	Last    int      `json:"last"`
	Metrics []Metric `json:"metrics"`
	Filter  string   `json:"filter"`
}

type Metric struct {
	Id           string      `json:"id"`
	Aggregations Aggregation `json:"aggregations"`
}

type Aggregation struct {
	Time  string `json:"time"`
	Group string `json:"group"`
}

type Period struct {
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type ApiResult struct {
	Data  []Data `json:"data"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type Data struct {
	Value []int `json:"d"`
}

func (s *Sysdigclient) GetSumMetric(metrics []Metric, filter string, period Period) (int, error) {

	secondsPeriod, err := s.getPeriodInSeconds(period)

	if err != nil {
		return 0, fmt.Errorf("error on give a period in seconds: %s", err)
	}

	query := Query{
		Last:    secondsPeriod,
		Filter:  filter,
		Metrics: metrics,
	}

	bodyValue, err := json.Marshal(query)

	if err != nil {
		return 0, fmt.Errorf("error on marshal query: %s", err)
	}

	response := s.client.DoRequest(
		client.Request{
			Method: "POST",
			URI:    "/api/data",
			Body:   bodyValue,
		},
	)

	var result ApiResult
	err = json.Unmarshal(response.Body, &result)

	if err != nil {
		return 0, fmt.Errorf("error on unmarshal response result: %s", err)
	}

	return result.Data[0].Value[0], nil
}

func (s *Sysdigclient) getPeriodInSeconds(period Period) (int, error) {
	if period.Days > 0 {
		return period.Days * 24 * 60 * 60, nil
	} else if period.Hours > 0 {
		return period.Hours * 60 * 60, nil
	} else if period.Minutes > 0 {
		return period.Minutes * 60, nil
	} else {
		return 0, fmt.Errorf("one of the period fields must be greater than zero")
	}
}

func New() *Sysdigclient {
	return &Sysdigclient{
		client: client.Client{
			URL: "https://app.sysdigcloud.com",
		},
	}
}

func NewWithUrl(url string) *Sysdigclient {
	return &Sysdigclient{
		client: client.Client{
			URL: url,
		},
	}
}
