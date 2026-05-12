package cta

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type GetStopsResponse struct {
	BustimeResponse struct {
		Stops []Stop `json:"stops"`
	} `json:"bustime-response"`
}

type Stop struct {
	Stpid string  `json:"stpid"`
	Stpnm string  `json:"stpnm"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
}

func (c *Client) GetStops(rt, dir string) (*GetStopsResponse, error) {
	u, err := url.Parse(bustimeV3Base + "/getstops")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("format", "json")
	q.Set("rt", rt)
	q.Set("dir", dir)
	u.RawQuery = q.Encode()

	start := time.Now()
	resp, err := http.Get(u.String())
	dur := time.Since(start)
	if err != nil {
		log.Printf("cta api: getstops rt=%s dir=%s status=- duration=%v err=%v", rt, dir, dur, err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("cta api: getstops rt=%s dir=%s status=%d duration=%v", rt, dir, resp.StatusCode, dur)

	var out GetStopsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
