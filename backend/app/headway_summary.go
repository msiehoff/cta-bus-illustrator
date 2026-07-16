package app

import (
	"math"
	"sort"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// SummarizeMinutes computes mean, median, sample stddev, CV, and avg wait from gap lengths.
func SummarizeMinutes(minutes []float64) business.HeadwaySummaryStats {
	n := len(minutes)
	if n == 0 {
		return business.HeadwaySummaryStats{}
	}

	sorted := append([]float64(nil), minutes...)
	sort.Float64s(sorted)

	var sum float64
	for _, m := range sorted {
		sum += m
	}
	mean := sum / float64(n)
	median := medianSorted(sorted)

	var stddev float64
	if n >= 2 {
		var sumSq float64
		for _, m := range sorted {
			d := m - mean
			sumSq += d * d
		}
		stddev = math.Sqrt(sumSq / float64(n-1))
	}

	var cv float64
	if mean > 0 && n >= 2 {
		cv = stddev / mean
	}

	return business.HeadwaySummaryStats{
		Count:          n,
		MeanMinutes:    mean,
		MedianMinutes:  median,
		StdDevMinutes:  stddev,
		CV:             cv,
		AvgWaitMinutes: mean / 2,
	}
}

func medianSorted(sorted []float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	mid := n / 2
	if n%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}

// SummarizeHeadways computes overall stats for a slice of observed headways.
func SummarizeHeadways(headways []business.Headway) business.HeadwaySummaryStats {
	mins := make([]float64, len(headways))
	for i, h := range headways {
		mins[i] = h.HeadwayMinutes
	}
	return SummarizeMinutes(mins)
}

// SummarizeHeadwaysByStop groups by stop/route/direction and summarizes each group.
// When aggregating to route+direction, prefer averaging per-stop means (equal stop weight).
func SummarizeHeadwaysByStop(headways []business.Headway) []business.HeadwayStopSummary {
	type key struct {
		StopID, RouteID, Direction string
	}
	groups := make(map[key][]float64)
	meta := make(map[key]business.Headway)
	order := make([]key, 0)

	for _, h := range headways {
		k := key{StopID: h.StopID, RouteID: h.RouteID, Direction: h.Direction}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
			meta[k] = h
		}
		groups[k] = append(groups[k], h.HeadwayMinutes)
	}

	out := make([]business.HeadwayStopSummary, 0, len(order))
	for _, k := range order {
		h := meta[k]
		out = append(out, business.HeadwayStopSummary{
			StopID:              k.StopID,
			StopName:            h.StopName,
			RouteID:             k.RouteID,
			RouteName:           h.RouteName,
			Direction:           k.Direction,
			HeadwaySummaryStats: SummarizeMinutes(groups[k]),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].MeanMinutes != out[j].MeanMinutes {
			return out[i].MeanMinutes > out[j].MeanMinutes
		}
		return out[i].StopName < out[j].StopName
	})
	return out
}

// MeanOfStopMeans averages per-stop mean headways (equal weight per stop).
// Count is the sum of per-stop counts. Median/CV/stddev are means of the per-stop metrics.
func MeanOfStopMeans(stops []business.HeadwayStopSummary) business.HeadwaySummaryStats {
	if len(stops) == 0 {
		return business.HeadwaySummaryStats{}
	}

	var (
		totalCount int
		sumMean    float64
		sumMedian  float64
		sumStd     float64
		sumCV      float64
		nStops     float64
	)
	for _, s := range stops {
		if s.Count == 0 {
			continue
		}
		totalCount += s.Count
		sumMean += s.MeanMinutes
		sumMedian += s.MedianMinutes
		sumStd += s.StdDevMinutes
		sumCV += s.CV
		nStops++
	}
	if nStops == 0 {
		return business.HeadwaySummaryStats{}
	}

	mean := sumMean / nStops
	return business.HeadwaySummaryStats{
		Count:          totalCount,
		MeanMinutes:    mean,
		MedianMinutes:  sumMedian / nStops,
		StdDevMinutes:  sumStd / nStops,
		CV:             sumCV / nStops,
		AvgWaitMinutes: mean / 2,
	}
}

// BuildPersistedSummaries creates stop, route_direction, route, and service_day
// summary rows for a completed service date from observed headway gaps.
func BuildPersistedSummaries(
	headways []business.Headway,
	serviceDate time.Time,
	windowStart, windowEnd time.Time,
) []business.HeadwaySummary {
	if len(headways) == 0 {
		return nil
	}

	byStop := SummarizeHeadwaysByStop(headways)
	out := make([]business.HeadwaySummary, 0, len(byStop)*2+16)

	for _, s := range byStop {
		out = append(out, business.HeadwaySummary{
			ServiceDate:         serviceDate,
			WindowStart:         windowStart,
			WindowEnd:           windowEnd,
			Grain:               business.HeadwayGrainStop,
			Method:              business.HeadwayMethodPooled,
			StopID:              s.StopID,
			RouteID:             s.RouteID,
			Direction:           s.Direction,
			HeadwaySummaryStats: s.HeadwaySummaryStats,
		})
	}

	type rdKey struct{ RouteID, Direction string }
	rdStops := make(map[rdKey][]business.HeadwayStopSummary)
	rdGaps := make(map[rdKey][]float64)
	rdOrder := make([]rdKey, 0)

	routeStops := make(map[string][]business.HeadwayStopSummary)
	routeGaps := make(map[string][]float64)
	routeOrder := make([]string, 0)

	for _, h := range headways {
		k := rdKey{RouteID: h.RouteID, Direction: h.Direction}
		if _, ok := rdGaps[k]; !ok {
			rdOrder = append(rdOrder, k)
		}
		rdGaps[k] = append(rdGaps[k], h.HeadwayMinutes)

		if _, ok := routeGaps[h.RouteID]; !ok {
			routeOrder = append(routeOrder, h.RouteID)
		}
		routeGaps[h.RouteID] = append(routeGaps[h.RouteID], h.HeadwayMinutes)
	}
	for _, s := range byStop {
		k := rdKey{RouteID: s.RouteID, Direction: s.Direction}
		rdStops[k] = append(rdStops[k], s)
		routeStops[s.RouteID] = append(routeStops[s.RouteID], s)
	}

	for _, k := range rdOrder {
		out = append(out,
			business.HeadwaySummary{
				ServiceDate:         serviceDate,
				WindowStart:         windowStart,
				WindowEnd:           windowEnd,
				Grain:               business.HeadwayGrainRouteDirection,
				Method:              business.HeadwayMethodPooled,
				RouteID:             k.RouteID,
				Direction:           k.Direction,
				HeadwaySummaryStats: SummarizeMinutes(rdGaps[k]),
			},
			business.HeadwaySummary{
				ServiceDate:         serviceDate,
				WindowStart:         windowStart,
				WindowEnd:           windowEnd,
				Grain:               business.HeadwayGrainRouteDirection,
				Method:              business.HeadwayMethodEqualStop,
				RouteID:             k.RouteID,
				Direction:           k.Direction,
				HeadwaySummaryStats: MeanOfStopMeans(rdStops[k]),
			},
		)
	}

	for _, routeID := range routeOrder {
		out = append(out,
			business.HeadwaySummary{
				ServiceDate:         serviceDate,
				WindowStart:         windowStart,
				WindowEnd:           windowEnd,
				Grain:               business.HeadwayGrainRoute,
				Method:              business.HeadwayMethodPooled,
				RouteID:             routeID,
				HeadwaySummaryStats: SummarizeMinutes(routeGaps[routeID]),
			},
			business.HeadwaySummary{
				ServiceDate:         serviceDate,
				WindowStart:         windowStart,
				WindowEnd:           windowEnd,
				Grain:               business.HeadwayGrainRoute,
				Method:              business.HeadwayMethodEqualStop,
				RouteID:             routeID,
				HeadwaySummaryStats: MeanOfStopMeans(routeStops[routeID]),
			},
		)
	}

	out = append(out,
		business.HeadwaySummary{
			ServiceDate:         serviceDate,
			WindowStart:         windowStart,
			WindowEnd:           windowEnd,
			Grain:               business.HeadwayGrainServiceDay,
			Method:              business.HeadwayMethodPooled,
			HeadwaySummaryStats: SummarizeHeadways(headways),
		},
		business.HeadwaySummary{
			ServiceDate:         serviceDate,
			WindowStart:         windowStart,
			WindowEnd:           windowEnd,
			Grain:               business.HeadwayGrainServiceDay,
			Method:              business.HeadwayMethodEqualStop,
			HeadwaySummaryStats: MeanOfStopMeans(byStop),
		},
	)

	return out
}
