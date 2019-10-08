// +build unit

package client_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NeowayLabs/sysdig-client/client"
)

func TestDoRequestExecuteReturnCorrectly(t *testing.T) {
	methods := []string{"GET", "POST"}
	for _, m := range methods {
		body := json.RawMessage(`{"test":1}`)
		location := "/newlocation"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == m {
				w.Header().Set("Location", location)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(body)
			}
		}))

		c := client.New(server.URL)

		actual := c.DoRequest(
			client.Request{
				Method: m,
				URI:    "/teste",
				Body:   body,
			},
		)

		expected := client.Response{
			Body:     body,
			Status:   http.StatusOK,
			Location: location,
			Error:    nil,
		}

		assertResponse(actual, expected, t)
	}
}

func TestDoRequestExecuteReturnsError(t *testing.T) {
	methods := []string{"GET", "POST"}
	for _, m := range methods {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == m {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		c := client.New(server.URL)
		actual := c.DoRequest(
			client.Request{
				Method: m,
				URI:    "/teste",
				Body:   nil,
			},
		)

		expected := client.Response{
			Body:     []byte(""),
			Status:   http.StatusInternalServerError,
			Location: "",
			Error:    nil,
		}

		assertResponse(actual, expected, t)
	}
}

func TestDoRequestExecuteReturnsErrorWhenUnauthorized(t *testing.T) {
	methods := []string{"GET", "POST"}
	for _, m := range methods {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == m {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}))

		c := client.New(server.URL)
		actual := c.DoRequest(
			client.Request{
				Method: m,
				URI:    "/teste",
				Body:   nil,
			},
		)

		assert.Error(t, actual.Error, "invalid access token, please enter correct key in environment variable SYSDIG_CLOUD_API_TOKEN")
		assert.Equal(t, actual.Status, http.StatusUnauthorized)
	}
}

func TestDoRequestExecuteReturnsErrorWhenBadRequestCredentials(t *testing.T) {
	result := map[string]interface{}{
		"timestamp": 1570537792986,
		"status":    400,
		"error":     "Bad Request",
		"message":   "Following header must be provided: X-Sysdig-Product",
		"path":      "/api/data",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/teste" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result)
		}
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	actual := c.DoRequest(
		client.Request{
			Method: "POST",
			URI:    "/teste",
			Body:   nil,
		},
	)

	assert.EqualError(t, actual.Error, "please set the variable SYSDIG_CLOUD_API_TOKEN with the token with the pattern `Bearer your_token`")
	assert.Equal(t, http.StatusUnauthorized, actual.Status)
}

func TestDoRequestExecuteReturnsErrorWhenBadRequestCredentialsAndUnmarshalError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/teste" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(2)
		}
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	actual := c.DoRequest(
		client.Request{
			Method: "POST",
			URI:    "/teste",
			Body:   nil,
		},
	)

	assert.EqualError(t, actual.Error, "error on unmarshal response body")
}

func assertResponse(actual, expected client.Response, t *testing.T) {
	if actual.Status != expected.Status {
		t.Fatal("Status should be equal! Actual: ", actual.Status, "Expected: ", expected.Status)
	}
	if strings.Trim(string(actual.Body), "\n") != string(expected.Body) {
		t.Fatal("Response should be equal! Actual: ", strings.Trim(string(actual.Body), "\n"), "Expected: ", string(expected.Body))
	}
	if actual.Location != expected.Location {
		t.Fatal("Location should be equal! Actual: ", actual.Location, "Expected: ", expected.Location)
	}
	if actual.Error != expected.Error {
		t.Fatal("Error should be equal! Actual: ", actual.Error, "Expected: ", expected.Error)
	}
}
