package api

import (
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type AdminSessionResponse struct {
	Authenticated bool   `json:"authenticated"`
	Username      string `json:"username,omitempty"`
}

type PipelineStatusResponse struct {
	Enabled       bool     `json:"enabled"`
	Running       bool     `json:"running"`
	Routes        []string `json:"routes"`
	RouteCount    int      `json:"routeCount"`
	PollInterval  string   `json:"pollInterval"`
	LastPollAt    *string  `json:"lastPollAt,omitempty"`
	LastPingCount int      `json:"lastPingCount"`
	LastError     string   `json:"lastError,omitempty"`
	StartedAt     *string  `json:"startedAt,omitempty"`
	ArrivalCount  int64    `json:"arrivalCount"`
}

type ArrivalResponse struct {
	StopID    string `json:"stopId"`
	StopName  string `json:"stopName,omitempty"`
	RouteID   string `json:"routeId"`
	RouteName string `json:"routeName,omitempty"`
	Direction string `json:"direction"`
	VehicleID string `json:"vehicleId"`
	Timestamp string `json:"timestamp"`
}

type ListArrivalsResponse struct {
	Arrivals []ArrivalResponse `json:"arrivals"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type HeadwayJobRunResponse struct {
	ID                int64   `json:"id"`
	ServiceDate       string  `json:"serviceDate"`
	Status            string  `json:"status"`
	TriggeredBy       string  `json:"triggeredBy"`
	StartedAt         string  `json:"startedAt"`
	FinishedAt        *string `json:"finishedAt,omitempty"`
	ArrivalsProcessed int     `json:"arrivalsProcessed"`
	HeadwaysWritten   int     `json:"headwaysWritten"`
	SummariesWritten  int     `json:"summariesWritten"`
	ErrorMessage      string  `json:"errorMessage,omitempty"`
}

type ListHeadwayJobRunsResponse struct {
	Runs   []HeadwayJobRunResponse `json:"runs"`
	Limit  int                     `json:"limit"`
	Offset int                     `json:"offset"`
}

type HeadwayResponse struct {
	StopID         string  `json:"stopId"`
	StopName       string  `json:"stopName,omitempty"`
	RouteID        string  `json:"routeId"`
	RouteName      string  `json:"routeName,omitempty"`
	Direction      string  `json:"direction"`
	Timestamp      string  `json:"timestamp"`
	HeadwayMinutes float64 `json:"headwayMinutes"`
	FromVehicleID  string  `json:"fromVehicleId,omitempty"`
	ToVehicleID    string  `json:"toVehicleId,omitempty"`
}

type ListHeadwaysResponse struct {
	Headways []HeadwayResponse `json:"headways"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type HeadwaySummaryStatsResponse struct {
	Count          int     `json:"count"`
	MeanMinutes    float64 `json:"meanMinutes"`
	MedianMinutes  float64 `json:"medianMinutes"`
	StdDevMinutes  float64 `json:"stdDevMinutes"`
	CV             float64 `json:"cv"`
	AvgWaitMinutes float64 `json:"avgWaitMinutes"`
}

type HeadwayStopSummaryResponse struct {
	StopID    string `json:"stopId"`
	StopName  string `json:"stopName,omitempty"`
	RouteID   string `json:"routeId"`
	RouteName string `json:"routeName,omitempty"`
	Direction string `json:"direction"`
	HeadwaySummaryStatsResponse
}

type HeadwaySummaryResponse struct {
	// Pooled: stats over all matching observed gaps (weighted by volume).
	Pooled HeadwaySummaryStatsResponse `json:"pooled"`
	// EqualStopWeight: mean of per-stop means (and related) when multiple stops.
	EqualStopWeight HeadwaySummaryStatsResponse  `json:"equalStopWeight"`
	ByStop          []HeadwayStopSummaryResponse `json:"byStop"`
	// Source is "stored" when read from headway_summaries, else "computed".
	Source string `json:"source"`
}

func toSummaryStatsResponse(s business.HeadwaySummaryStats) HeadwaySummaryStatsResponse {
	return HeadwaySummaryStatsResponse{
		Count:          s.Count,
		MeanMinutes:    round2(s.MeanMinutes),
		MedianMinutes:  round2(s.MedianMinutes),
		StdDevMinutes:  round2(s.StdDevMinutes),
		CV:             round3(s.CV),
		AvgWaitMinutes: round2(s.AvgWaitMinutes),
	}
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func round3(v float64) float64 {
	return float64(int(v*1000+0.5)) / 1000
}

func HeadwayJobRunResponseFrom(run business.HeadwayJobRun) HeadwayJobRunResponse {
	resp := HeadwayJobRunResponse{
		ID:                run.ID,
		ServiceDate:       run.ServiceDate.UTC().Format("2006-01-02"),
		Status:            string(run.Status),
		TriggeredBy:       string(run.TriggeredBy),
		StartedAt:         run.StartedAt.UTC().Format(time.RFC3339),
		ArrivalsProcessed: run.ArrivalsProcessed,
		HeadwaysWritten:   run.HeadwaysWritten,
		SummariesWritten:  run.SummariesWritten,
		ErrorMessage:      run.ErrorMessage,
	}
	if run.FinishedAt != nil {
		v := run.FinishedAt.UTC().Format(time.RFC3339)
		resp.FinishedAt = &v
	}
	return resp
}

func toPipelineStatusResponse(status app.PipelineStatus, enabled bool, arrivalCount int64) PipelineStatusResponse {
	resp := PipelineStatusResponse{
		Enabled:       enabled,
		Running:       status.Running,
		Routes:        status.Routes,
		RouteCount:    len(status.Routes),
		PollInterval:  status.PollInterval.String(),
		LastPingCount: status.LastPingCount,
		LastError:     status.LastError,
		ArrivalCount:  arrivalCount,
	}
	if !status.LastPollAt.IsZero() {
		v := status.LastPollAt.UTC().Format(time.RFC3339)
		resp.LastPollAt = &v
	}
	if !status.StartedAt.IsZero() {
		v := status.StartedAt.UTC().Format(time.RFC3339)
		resp.StartedAt = &v
	}
	return resp
}

func toArrivalResponse(arrival business.Arrival) ArrivalResponse {
	return ArrivalResponse{
		StopID:    arrival.StopID,
		StopName:  arrival.StopName,
		RouteID:   arrival.RouteID,
		RouteName: arrival.RouteName,
		Direction: arrival.Direction,
		VehicleID: arrival.VehicleID,
		Timestamp: arrival.Timestamp.UTC().Format(time.RFC3339),
	}
}

func toHeadwayResponse(h business.Headway) HeadwayResponse {
	return HeadwayResponse{
		StopID:         h.StopID,
		StopName:       h.StopName,
		RouteID:        h.RouteID,
		RouteName:      h.RouteName,
		Direction:      h.Direction,
		Timestamp:      h.Timestamp.UTC().Format(time.RFC3339),
		HeadwayMinutes: h.HeadwayMinutes,
		FromVehicleID:  h.FromVehicleID,
		ToVehicleID:    h.ToVehicleID,
	}
}
