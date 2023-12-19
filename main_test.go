package main

import (
	"go-short-url/cache"
	"go-short-url/config"
	"go-short-url/controller"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {

	tt := []struct {
		name   string
		status int
		method string
		err    string
	}{
		{name: "ok request", status: http.StatusOK, method: "GET", err: ""},
		{name: "forbidden method post", status: http.StatusMethodNotAllowed, method: "POST", err: "Forbidden HTTP Method"},
		{name: "forbidden method put", status: http.StatusMethodNotAllowed, method: "PUT", err: "Forbidden HTTP Method"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, config.HostURL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rec := httptest.NewRecorder()

			homeHandler(rec, req)

			res := rec.Result()
			// Test response is correct: 200 or 405
			if res.StatusCode != tc.status {
				t.Errorf("expected status %d; got %v", tc.status, res.StatusCode)
			}
			defer res.Body.Close()

			if tc.err != "" {
				if res.StatusCode != tc.status {
					t.Errorf("expected status %d; got: %+v", tc.status, res.StatusCode)
				}
			}

			_, err = io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %+v", err)
			}
		})
	}
}

func TestCreateURL(t *testing.T) {

	testURL := url.Values{
		"full_url": []string{
			"https://rshipp.com/go-api-integration-testing/",
		},
	}

	formData := strings.NewReader(testURL.Encode())

	req, err := http.NewRequest("POST", config.HostURL+"/createURL", formData)
	if err != nil {
		t.Fatalf("could not create request: %+v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	rec := httptest.NewRecorder()

	newCache := cache.NewCache()
	userURL := controller.NewUserURL(newCache)

	userURL.CreateURL(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got: %+v", http.StatusOK, res.StatusCode)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("could not read response: %+v", err)
	}

	// check len of response body
	if len(b) < 3000 && len(b) > 5000 {
		t.Errorf("expected response body lenght: %d; got: %d", 3000, len(b))
	}
}
