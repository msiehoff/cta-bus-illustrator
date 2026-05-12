package cta

import "strings"

const defaultBustimeV3Base = "https://www.ctabustracker.com/bustime/api/v3"

type Client struct {
	apiKey string
	// BaseURL is the bustime API v3 root (e.g. https://www.ctabustracker.com/bustime/api/v3).
	// If empty, defaultBustimeV3Base is used. Set to an httptest server URL in tests.
	BaseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) bustimeV3Base() string {
	s := c.BaseURL
	if s == "" {
		s = defaultBustimeV3Base
	}
	return strings.TrimSuffix(s, "/")
}
