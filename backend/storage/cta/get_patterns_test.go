package cta

import "testing"

type patternIn struct {
	rtdir string
	ln    float64
	lat   float64
}

func patternResp(entries ...patternIn) *GetRoutePatternResponse {
	var resp GetRoutePatternResponse
	for i, e := range entries {
		pt := struct {
			Seq   int     `json:"seq"`
			Lat   float64 `json:"lat"`
			Lon   float64 `json:"lon"`
			Typ   string  `json:"typ"`
			Stpid string  `json:"stpid,omitempty"`
			Stpnm string  `json:"stpnm,omitempty"`
			Pdist float64 `json:"pdist"`
		}{Seq: 1, Lat: e.lat, Lon: -87.6, Typ: "W"}
		resp.BustimeResponse.Ptr = append(resp.BustimeResponse.Ptr, struct {
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
		}{
			Pid:   i + 1,
			Ln:    e.ln,
			Rtdir: e.rtdir,
			Pt:    []struct {
				Seq   int     `json:"seq"`
				Lat   float64 `json:"lat"`
				Lon   float64 `json:"lon"`
				Typ   string  `json:"typ"`
				Stpid string  `json:"stpid,omitempty"`
				Stpnm string  `json:"stpnm,omitempty"`
				Pdist float64 `json:"pdist"`
			}{pt},
		})
	}
	return &resp
}

func TestSegmentsFromPatternResponsePrefersNorthbound(t *testing.T) {
	resp := patternResp(
		patternIn{"Southbound", 99999, 41.1},
		patternIn{"Northbound", 100, 41.9},
	)
	segs, err := SegmentsFromPatternResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(segs) != 1 || segs[0].Lat != 41.9 {
		t.Fatalf("expected Northbound point, got %+v", segs)
	}
}

func TestSegmentsFromPatternResponseFallsBackToSouthboundOnly(t *testing.T) {
	// CTA currently publishes only Southbound for some routes (e.g. 10, 171).
	resp := patternResp(patternIn{"Southbound", 54184, 41.75})
	segs, err := SegmentsFromPatternResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(segs) != 1 || segs[0].Lat != 41.75 {
		t.Fatalf("expected Southbound fallback, got %+v", segs)
	}
}

func TestSegmentsFromPatternResponseNormalizesPreferredDirection(t *testing.T) {
	resp := patternResp(
		patternIn{"WEST", 500, 41.2},
		patternIn{"East Bound", 200, 41.8},
	)
	segs, err := SegmentsFromPatternResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(segs) != 1 || segs[0].Lat != 41.8 {
		t.Fatalf("expected normalized Eastbound, got %+v", segs)
	}
}

func TestSegmentsFromPatternResponseEmpty(t *testing.T) {
	_, err := SegmentsFromPatternResponse(&GetRoutePatternResponse{})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}
