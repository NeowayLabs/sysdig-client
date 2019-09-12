package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const RequestTimeoutSeconds = 360
const timeout = time.Duration(time.Duration(RequestTimeoutSeconds) * time.Second)

type Client struct {
	URL        string
	HttpClient http.Client
}

type Request struct {
	Method string
	URI    string
	Body   json.RawMessage
}

type Response struct {
	Body       json.RawMessage
	Location   string
	RetryAfter int
	Status     int
	Error      error
}

func (c *Client) DoRequest(r Request) Response {
	response := Response{}

	b, _ := json.Marshal(r.Body)
	req, err := http.NewRequest(r.Method, c.URL+r.URI, bytes.NewBuffer(b))
	if err != nil {
		response.Error = err
		return response
	}

	req.Header.Set("Authorization", os.Getenv("SYSDIG_CLOUD_API_TOKEN"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		response.Error = err
		return response
	}
	defer resp.Body.Close()

	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
		return response
	}

	response.Body = msg
	response.Status = resp.StatusCode
	response.Location = resp.Header.Get("Location")
	value := resp.Header.Get("Retry-After")
	if value != "" {
		v, err := strconv.Atoi(value)
		if err != nil {
			response.Error = err
			return response
		}
		response.RetryAfter = v
	}
	response.Error = err

	return response
}

func New(url string) *Client {
	return &Client{
		URL:        url,
		HttpClient: http.Client{Timeout: timeout},
	}
}
