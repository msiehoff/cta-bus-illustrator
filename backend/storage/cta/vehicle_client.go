package cta

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// bustimeTimestampLayout is the CTA Bus Tracker timestamp format (America/Chicago).
const bustimeTimestampLayout = "20060102 15:04:05"

// VehicleClient adapts the low-level CTA HTTP client to the pipeline's CTAVehicleClient port.
type VehicleClient struct {
	client *Client
	loc    *time.Location
}

func NewVehicleClient(client *Client) *VehicleClient {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		loc = time.FixedZone("CST", -6*60*60)
	}
	return &VehicleClient{client: client, loc: loc}
}

func (v *VehicleClient) GetVehicles(ctx context.Context, routeIDs []string) ([]business.VehiclePing, error) {
	if len(routeIDs) == 0 {
		return nil, nil
	}
	_ = ctx

	resp, err := v.client.GetVehicles(strings.Join(routeIDs, ","))
	if err != nil {
		return nil, err
	}

	pings := make([]business.VehiclePing, 0, len(resp.BustimeResponse.Vehicle))
	for _, vehicle := range resp.BustimeResponse.Vehicle {
		ping, err := v.toVehiclePing(vehicle)
		if err != nil {
			continue
		}
		pings = append(pings, ping)
	}
	return pings, nil
}

func (v *VehicleClient) GetStops(ctx context.Context, routeID, direction string) ([]business.Stop, error) {
	_ = ctx

	resp, err := v.client.GetStops(routeID, direction)
	if err != nil {
		return nil, err
	}

	stops := make([]business.Stop, len(resp.BustimeResponse.Stops))
	for i, s := range resp.BustimeResponse.Stops {
		stops[i] = business.Stop{
			StopID:    s.Stpid,
			RouteID:   routeID,
			Direction: direction,
			Name:      s.Stpnm,
			Lat:       s.Lat,
			Lon:       s.Lon,
			Sequence:  i + 1,
		}
	}
	return stops, nil
}

func (v *VehicleClient) toVehiclePing(vehicle Vehicle) (business.VehiclePing, error) {
	lat, err := strconv.ParseFloat(vehicle.Lat, 64)
	if err != nil {
		return business.VehiclePing{}, fmt.Errorf("parse lat for vid %s: %w", vehicle.Vid, err)
	}
	lon, err := strconv.ParseFloat(vehicle.Lon, 64)
	if err != nil {
		return business.VehiclePing{}, fmt.Errorf("parse lon for vid %s: %w", vehicle.Vid, err)
	}

	ts, err := time.ParseInLocation(bustimeTimestampLayout, vehicle.Tmstmp, v.loc)
	if err != nil {
		ts = time.Now()
	}

	direction := strings.TrimSpace(vehicle.Stsd)
	if direction == "" {
		direction = strings.TrimSpace(vehicle.Des)
	}

	return business.VehiclePing{
		VehicleID: vehicle.Vid,
		RouteID:   vehicle.Rt,
		Direction: direction,
		Lat:       lat,
		Lon:       lon,
		Timestamp: ts,
	}, nil
}
