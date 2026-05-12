package cta

const bustimeV3Base = "https://www.ctabustracker.com/bustime/api/v3"

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
