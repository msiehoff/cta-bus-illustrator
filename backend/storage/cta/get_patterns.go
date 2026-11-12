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
	url := fmt.Sprintf("%s/getpatterns?key=%s&format=json&rt=%s", c.bustimeV3Base(), c.apiKey, routeID)
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

func isPreferredMapDirection(rtdir string) bool {
	switch business.NormalizeDirection(rtdir) {
	case "Northbound", "Eastbound":
		return true
	default:
		return false
	}
}

// SegmentsFromPatternResponse maps a getpatterns JSON body to route segments.
// Prefer Northbound/Eastbound patterns (skip Southbound/Westbound when a preferred
// direction exists) so the map draws one representative path. When CTA only publishes
// Southbound/Westbound (or another direction) for a route — e.g. rt 10, 171 — fall back
// to the longest available pattern so import-segments can succeed.
// When multiple patterns exist in the chosen set (alternate routings), pick the single
// pattern with the largest ln (CTA-reported path length), not all of them: concatenating
// multiple polylines would connect the end of one path to the start of another and draw
// long bogus diagonals on the map.
func SegmentsFromPatternResponse(resp *GetRoutePatternResponse) ([]business.RouteSegment, error) {
	if resp == nil || len(resp.BustimeResponse.Ptr) == 0 {
		return nil, fmt.Errorf("no pattern directions in response")
	}

	bestPreferredIdx := -1
	bestPreferredLen := -1.0
	bestAnyIdx := -1
	bestAnyLen := -1.0
	seen := make([]string, 0, len(resp.BustimeResponse.Ptr))

	for i := range resp.BustimeResponse.Ptr {
		d := resp.BustimeResponse.Ptr[i]
		seen = append(seen, d.Rtdir)
		if d.Ln > bestAnyLen {
			bestAnyLen = d.Ln
			bestAnyIdx = i
		}
		if !isPreferredMapDirection(d.Rtdir) {
			continue
		}
		if d.Ln > bestPreferredLen {
			bestPreferredLen = d.Ln
			bestPreferredIdx = i
		}
	}

	bestIdx := bestPreferredIdx
	if bestIdx < 0 {
		bestIdx = bestAnyIdx
	}
	if bestIdx < 0 {
		return nil, fmt.Errorf("no usable pattern directions in response (got %q)", seen)
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
