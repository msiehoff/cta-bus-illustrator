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
	RouteID   string `json:"routeId"`
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
		RouteID:   arrival.RouteID,
		Direction: arrival.Direction,
		VehicleID: arrival.VehicleID,
		Timestamp: arrival.Timestamp.UTC().Format(time.RFC3339),
	}
}
