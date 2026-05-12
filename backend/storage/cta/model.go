package cta

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

func (c *Client) GetRoutePattern(routeID string) (*GetRoutePatternResponse, error) {
	url := fmt.Sprintf("%s/getpatterns?key=%s&format=json&rt=%s", bustimeV3Base, c.apiKey, routeID)
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

// SegmentsFromPatternResponse maps a getpatterns JSON body to route segments.
// It considers only Northbound/Eastbound ptr entries (skips Southbound/Westbound).
// When multiple such patterns exist (alternate routings, e.g. State vs Halsted),
// it picks the single pattern with the largest ln (CTA-reported path length), not
// all of them: concatenating multiple polylines would connect the end of one path
// to the start of another and draw long bogus diagonals on the map.
func SegmentsFromPatternResponse(resp *GetRoutePatternResponse) ([]business.RouteSegment, error) {
	if resp == nil || len(resp.BustimeResponse.Ptr) == 0 {
		return nil, fmt.Errorf("no pattern directions in response")
	}

	bestIdx := -1
	bestLen := -1.0
	for i := range resp.BustimeResponse.Ptr {
		d := resp.BustimeResponse.Ptr[i]
		if d.Rtdir != "Northbound" && d.Rtdir != "Eastbound" {
			continue
		}
		if d.Ln > bestLen {
			bestLen = d.Ln
			bestIdx = i
		}
	}

	if bestIdx < 0 {
		return nil, fmt.Errorf("no northbound/eastbound pattern directions in response")
	}

	direction := resp.BustimeResponse.Ptr[bestIdx]
	out := make([]business.RouteSegment, 0, len(direction.Pt))
	for _, p := range direction.Pt {
		out = append(out, business.RouteSegment{
			Sequence: p.Seq,
			Lat:      p.Lat,
			Lng:      p.Lon,
		})
	}
	return out, nil
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
