package app

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

const (
	DefaultHeadwayPeriodDays = 30
	MaxHeadwayPeriodDays     = 90
)

// HeadwayPeriodStats is an observation-weighted rollup of daily summary rows.
type HeadwayPeriodStats struct {
	business.HeadwaySummaryStats
	DaysWithData int
	PeriodStart  time.Time
	PeriodEnd    time.Time
}

// HeadwayDayPoint is one day's route or network summary for a chart series.
type HeadwayDayPoint struct {
	ServiceDate    time.Time
	MedianMinutes  float64
	AvgWaitMinutes float64
	CV             float64
	Count          int
}

// HeadwayRoutePeriod is period stats for one route (both directions).
type HeadwayRoutePeriod struct {
	RouteID   string
	RouteName string
	HeadwayPeriodStats
}

// HeadwayRouteDetail is period stats plus a daily series for one route.
type HeadwayRouteDetail struct {
	HeadwayRoutePeriod
	Series []HeadwayDayPoint
}

// HeadwaySystemOverview is network period stats, daily series, and shortest-headway routes.
type HeadwaySystemOverview struct {
	HeadwayPeriodStats
	Series          []HeadwayDayPoint
	ShortestHeadways []HeadwayRoutePeriod
}

// HeadwayPublicService reads persisted summaries for the rider-facing UI.
type HeadwayPublicService struct {
	summaries HeadwaySummaryRepository
}

func NewHeadwayPublicService(summaries HeadwaySummaryRepository) *HeadwayPublicService {
	return &HeadwayPublicService{summaries: summaries}
}

// ResolvePeriodDays clamps the requested day count.
func ResolvePeriodDays(days int) int {
	if days <= 0 {
		return DefaultHeadwayPeriodDays
	}
	if days > MaxHeadwayPeriodDays {
		return MaxHeadwayPeriodDays
	}
	return days
}

// ListRoutes returns one period row per route (grain=route, equal_stop).
func (s *HeadwayPublicService) ListRoutes(ctx context.Context, days int) ([]HeadwayRoutePeriod, HeadwayPeriodStats, error) {
	days = ResolvePeriodDays(days)
	rows, periodMeta, err := s.loadGrainRows(ctx, business.HeadwayGrainRoute, "", days)
	if err != nil {
		return nil, HeadwayPeriodStats{}, err
	}
	if len(rows) == 0 {
		return []HeadwayRoutePeriod{}, periodMeta, nil
	}

	byRoute := groupByRoute(rows)
	out := make([]HeadwayRoutePeriod, 0, len(byRoute))
	for routeID, routeRows := range byRoute {
		stats := AggregatePeriodStats(routeRows)
		stats.PeriodStart = periodMeta.PeriodStart
		stats.PeriodEnd = periodMeta.PeriodEnd
		out = append(out, HeadwayRoutePeriod{
			RouteID:            routeID,
			RouteName:          routeNameFrom(routeRows),
			HeadwayPeriodStats: stats,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].MedianMinutes != out[j].MedianMinutes {
			return out[i].MedianMinutes > out[j].MedianMinutes
		}
		return out[i].RouteID < out[j].RouteID
	})
	return out, periodMeta, nil
}

// GetRoute returns period stats and daily series for one route.
func (s *HeadwayPublicService) GetRoute(ctx context.Context, routeID string, days int) (HeadwayRouteDetail, error) {
	if routeID == "" {
		return HeadwayRouteDetail{}, fmt.Errorf("route id required")
	}
	days = ResolvePeriodDays(days)
	rows, periodMeta, err := s.loadGrainRows(ctx, business.HeadwayGrainRoute, routeID, days)
	if err != nil {
		return HeadwayRouteDetail{}, err
	}
	if len(rows) == 0 {
		return HeadwayRouteDetail{
			HeadwayRoutePeriod: HeadwayRoutePeriod{
				RouteID:            routeID,
				HeadwayPeriodStats: periodMeta,
			},
			Series: []HeadwayDayPoint{},
		}, nil
	}

	stats := AggregatePeriodStats(rows)
	stats.PeriodStart = periodMeta.PeriodStart
	stats.PeriodEnd = periodMeta.PeriodEnd
	return HeadwayRouteDetail{
		HeadwayRoutePeriod: HeadwayRoutePeriod{
			RouteID:            routeID,
			RouteName:          routeNameFrom(rows),
			HeadwayPeriodStats: stats,
		},
		Series: toDaySeries(rows),
	}, nil
}

// GetSystem returns network period stats, daily series, and top longest median routes.
func (s *HeadwayPublicService) GetSystem(ctx context.Context, days int) (HeadwaySystemOverview, error) {
	days = ResolvePeriodDays(days)

	dayRows, periodMeta, err := s.loadGrainRows(ctx, business.HeadwayGrainServiceDay, "", days)
	if err != nil {
		return HeadwaySystemOverview{}, err
	}

	routes, _, err := s.ListRoutes(ctx, days)
	if err != nil {
		return HeadwaySystemOverview{}, err
	}

	stats := AggregatePeriodStats(dayRows)
	stats.PeriodStart = periodMeta.PeriodStart
	stats.PeriodEnd = periodMeta.PeriodEnd

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].MedianMinutes != routes[j].MedianMinutes {
			return routes[i].MedianMinutes < routes[j].MedianMinutes
		}
		return routes[i].RouteID < routes[j].RouteID
	})
	shortest := routes
	if len(shortest) > 10 {
		shortest = shortest[:10]
	}

	return HeadwaySystemOverview{
		HeadwayPeriodStats: stats,
		Series:             toDaySeries(dayRows),
		ShortestHeadways:   shortest,
	}, nil
}

func (s *HeadwayPublicService) loadGrainRows(
	ctx context.Context,
	grain, routeID string,
	days int,
) ([]business.HeadwaySummary, HeadwayPeriodStats, error) {
	empty := HeadwayPeriodStats{}
	if s.summaries == nil {
		return nil, empty, nil
	}

	// Pull recent equal_stop rows; period is the N most recent distinct service dates.
	rows, err := s.summaries.List(ctx, HeadwaySummaryFilter{
		Grain:   grain,
		Method:  business.HeadwayMethodEqualStop,
		RouteID: routeID,
		SortAsc: false,
		Limit:   10_000,
	})
	if err != nil {
		return nil, empty, err
	}
	if len(rows) == 0 {
		return nil, empty, nil
	}

	dates := distinctServiceDates(rows)
	if len(dates) > days {
		dates = dates[:days]
	}
	allowed := make(map[string]struct{}, len(dates))
	for _, d := range dates {
		allowed[d.Format("2006-01-02")] = struct{}{}
	}

	filtered := make([]business.HeadwaySummary, 0, len(rows))
	for _, row := range rows {
		if _, ok := allowed[row.ServiceDate.Format("2006-01-02")]; ok {
			filtered = append(filtered, row)
		}
	}

	meta := HeadwayPeriodStats{
		DaysWithData: len(dates),
		PeriodEnd:    dates[0],
		PeriodStart:  dates[len(dates)-1],
	}
	return filtered, meta, nil
}

// AggregatePeriodStats observation-weights daily summary rows into one period snapshot.
func AggregatePeriodStats(rows []business.HeadwaySummary) HeadwayPeriodStats {
	if len(rows) == 0 {
		return HeadwayPeriodStats{}
	}

	dates := make(map[string]struct{})
	var (
		totalCount int
		sumMean    float64
		sumMedian  float64
		sumStd     float64
		sumCV      float64
		sumWait    float64
	)
	for _, row := range rows {
		if row.Count <= 0 {
			continue
		}
		n := float64(row.Count)
		totalCount += row.Count
		sumMean += row.MeanMinutes * n
		sumMedian += row.MedianMinutes * n
		sumStd += row.StdDevMinutes * n
		sumCV += row.CV * n
		sumWait += row.AvgWaitMinutes * n
		dates[row.ServiceDate.Format("2006-01-02")] = struct{}{}
	}
	if totalCount == 0 {
		return HeadwayPeriodStats{DaysWithData: len(dates)}
	}

	tn := float64(totalCount)
	mean := sumMean / tn
	return HeadwayPeriodStats{
		HeadwaySummaryStats: business.HeadwaySummaryStats{
			Count:          totalCount,
			MeanMinutes:    mean,
			MedianMinutes:  sumMedian / tn,
			StdDevMinutes:  sumStd / tn,
			CV:             sumCV / tn,
			AvgWaitMinutes: sumWait / tn,
		},
		DaysWithData: len(dates),
	}
}

func distinctServiceDates(rows []business.HeadwaySummary) []time.Time {
	seen := make(map[string]time.Time)
	for _, row := range rows {
		key := row.ServiceDate.Format("2006-01-02")
		if _, ok := seen[key]; ok {
			continue
		}
		day := time.Date(row.ServiceDate.Year(), row.ServiceDate.Month(), row.ServiceDate.Day(), 0, 0, 0, 0, time.UTC)
		seen[key] = day
	}
	out := make([]time.Time, 0, len(seen))
	for _, d := range seen {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].After(out[j]) })
	return out
}

func groupByRoute(rows []business.HeadwaySummary) map[string][]business.HeadwaySummary {
	out := make(map[string][]business.HeadwaySummary)
	for _, row := range rows {
		if row.RouteID == "" {
			continue
		}
		out[row.RouteID] = append(out[row.RouteID], row)
	}
	return out
}

func routeNameFrom(rows []business.HeadwaySummary) string {
	for _, row := range rows {
		if row.RouteName != "" {
			return row.RouteName
		}
	}
	return ""
}

func toDaySeries(rows []business.HeadwaySummary) []HeadwayDayPoint {
	byDay := make(map[string]business.HeadwaySummary)
	for _, row := range rows {
		key := row.ServiceDate.Format("2006-01-02")
		// Prefer keeping a single row per day (route/service_day grains are unique per day×method).
		if existing, ok := byDay[key]; ok {
			if row.Count > existing.Count {
				byDay[key] = row
			}
			continue
		}
		byDay[key] = row
	}

	keys := make([]string, 0, len(byDay))
	for k := range byDay {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]HeadwayDayPoint, 0, len(keys))
	for _, k := range keys {
		row := byDay[k]
		out = append(out, HeadwayDayPoint{
			ServiceDate:    row.ServiceDate,
			MedianMinutes:  row.MedianMinutes,
			AvgWaitMinutes: row.AvgWaitMinutes,
			CV:             row.CV,
			Count:          row.Count,
		})
	}
	return out
}
