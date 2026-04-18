package cta

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) GetRoutePattern(routeID string) (*GetRoutePatternResponse, error) {
	url := fmt.Sprintf("https://www.ctabustracker.com/bustime/api/v3/getpatterns?key=%s&format=json&rt=%s", c.apiKey, routeID)
	start := time.Now()
	resp, err := http.Get(url)
	dur := time.Since(start)
	if err != nil {
		log.Printf("cta api: getpatterns rt=%s status=- duration=%v err=%v", routeID, dur, err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("cta api: getpatterns rt=%s status=%d duration=%v", routeID, resp.StatusCode, dur)

	var response GetRoutePatternResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

type GetRoutePatternResponse struct {
	BustimeResponse struct {
		Ptr []struct {
			Pid   int     `json:"pid"`
			Ln    float64 `json:"ln"`
			Rtdir string  `json:"rtdir"`
			Pt    []struct {
				Seq   int     `json:"seq"`
				Lat   float64 `json:"lat"`
				Lon   float64 `json:"lon"`
				Typ   string  `json:"typ"`
				Stpid string  `json:"stpid,omitempty"`
				Stpnm string  `json:"stpnm,omitempty"`
				Pdist float64 `json:"pdist"`
			} `json:"pt"`
		} `json:"ptr"`
	} `json:"bustime-response"`
}
