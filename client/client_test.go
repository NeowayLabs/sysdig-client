// +build unit

package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NeowayLabs/sysdigcli/client"
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

		c := client.Client{}
		actual := c.DoRequest(
			client.Request{
				Method: m,
				URL:    server.URL,
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

		c := client.Client{}
		actual := c.DoRequest(
			client.Request{
				Method: m,
				URL:    server.URL,
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
