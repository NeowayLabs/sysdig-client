package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const RequestTimeoutSeconds = 360
const timeout = time.Duration(RequestTimeoutSeconds) * time.Second

type Client struct {
	Endpoint   string
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

	b, err := json.Marshal(r.Body)

	if err != nil {
		response.Error = err
		return response
	}

	req, err := http.NewRequest(r.Method, c.Endpoint+r.URI, bytes.NewBuffer(b))
	if err != nil {
		response.Error = err
		return response
	}

	authorization := fmt.Sprint(os.Getenv("SYSDIG_CLOUD_API_TOKEN"))
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		response.Error = err
		return response
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		response.Status = http.StatusUnauthorized
		response.Error = errors.New("unauthorized access token, please enter correct key in environment variable SYSDIG_CLOUD_API_TOKEN")
		return response
	} else if resp.StatusCode == http.StatusBadRequest {

		jsonValue := make(map[string]interface{})

		err = json.NewDecoder(resp.Body).Decode(&jsonValue)

		if err != nil {
			response.Error = errors.New("error on unmarshal response body")
			return response
		}

		if jsonValue["message"] == "Following header must be provided: X-Sysdig-Product" {
			response.Status = http.StatusUnauthorized
			response.Error = errors.New("please set the variable SYSDIG_CLOUD_API_TOKEN with the token with the pattern `Bearer your_token`")
			return response
		}
	}

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

func New(endpoint string) *Client {
	return &Client{
		Endpoint:   endpoint,
		HttpClient: http.Client{Timeout: timeout},
	}
}
