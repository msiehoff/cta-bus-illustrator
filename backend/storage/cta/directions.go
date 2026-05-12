package cta

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type GetDirectionsResponse struct {
	BustimeResponse struct {
		Directions []Direction `json:"directions"`
	} `json:"bustime-response"`
}

type Direction struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetDirections(rt string) (*GetDirectionsResponse, error) {
	u, err := url.Parse(bustimeV3Base + "/getdirections")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("format", "json")
	q.Set("rt", rt)
	u.RawQuery = q.Encode()

	start := time.Now()
	resp, err := http.Get(u.String())
	dur := time.Since(start)
	if err != nil {
		log.Printf("cta api: getdirections rt=%s status=- duration=%v err=%v", rt, dur, err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("cta api: getdirections rt=%s status=%d duration=%v", rt, resp.StatusCode, dur)

	var out GetDirectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
