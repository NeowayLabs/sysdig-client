package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const RequestTimeoutSeconds = 360

type Client struct {
	Authorization string
}

type Request struct {
	Method string
	URL    string
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
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(b))
	if err != nil {
		response.Error = err
		return response
	}

	timeout := time.Duration(time.Duration(RequestTimeoutSeconds) * time.Second)
	client := &http.Client{Timeout: timeout}

	req.Header.Set("Authorization", c.Authorization)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
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
