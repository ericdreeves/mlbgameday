// Package mlbgameday is a client library for the MLB Gameday API.
package mlbgameday

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client manages communication with the MLB Gameday API.
type Client struct {
	// HTTP client used to retrieve data from the MLB Gameday API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL
}

// NewClient returns a new MLB Gameday API client.
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	c := &Client{client: client}
	c.BaseURL = &url.URL{
		Scheme: "http",
		Host:   "gd2.mlb.com",
	}

	return c
}

// Gameday returns a new GamedayService for the provided date.
func (c *Client) Gameday(date time.Time) GamedayService {
	return NewGamedayService(c, date)
}

// get sends an HTTP GET to the MLB Gameday API at the requested path
// and returns the HTTP response body.
func (c *Client) get(path string) ([]byte, error) {
	rel, _ := url.Parse(path)
	u := c.BaseURL.ResolveReference(rel)

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %v", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}
