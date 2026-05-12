package cta

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type GetVehiclesResponse struct {
	BustimeResponse struct {
		Vehicle []Vehicle         `json:"vehicle"`
		Error   []BustimeAPIError `json:"error"`
	} `json:"bustime-response"`
}

type Vehicle struct {
	Vid          string `json:"vid"`
	Tmstmp       string `json:"tmstmp"`
	Lat          string `json:"lat"`
	Lon          string `json:"lon"`
	Hdg          string `json:"hdg"`
	Pid          int    `json:"pid"`
	Rt           string `json:"rt"`
	Des          string `json:"des"`
	Pdist        int    `json:"pdist"`
	Dly          bool   `json:"dly"`
	Tatripid     string `json:"tatripid"`
	Origtatripno string `json:"origtatripno"`
	Tablockid    string `json:"tablockid"`
	Zone         string `json:"zone"`
	Mode         int    `json:"mode"`
	Psgld        string `json:"psgld"`
	Stst         int    `json:"stst"`
	Stsd         string `json:"stsd"`
}

// BustimeAPIError is a route-level warning from the CTA API (e.g. unknown rt in a comma list).
type BustimeAPIError struct {
	Rt  string `json:"rt"`
	Msg string `json:"msg"`
}

func (c *Client) GetVehicles(rt string) (*GetVehiclesResponse, error) {
	u, err := url.Parse(c.bustimeV3Base() + "/getvehicles")
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
		log.Printf("cta api: getvehicles rt=%s status=- duration=%v err=%v", rt, dur, err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("cta api: getvehicles rt=%s status=%d duration=%v", rt, resp.StatusCode, dur)

	var out GetVehiclesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
