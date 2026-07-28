package app

import (
	"fmt"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// HeadwayWindowSpec describes a temporal slice of a service date (America/Chicago).
type HeadwayWindowSpec struct {
	Key   string
	Start time.Time // inclusive
	End   time.Time // exclusive
}

// ValidHeadwayWindowKeys are accepted by public/admin window= query params.
var ValidHeadwayWindowKeys = map[string]struct{}{
	business.HeadwayWindowOverall:    {},
	business.HeadwayWindowWeekdayAM:  {},
	business.HeadwayWindowWeekdayMid: {},
	business.HeadwayWindowWeekdayPM:  {},
	business.HeadwayWindowWeekdayEve: {},
	business.HeadwayWindowSaturday:   {},
	business.HeadwayWindowSunday:     {},
}

// ParseHeadwayWindowKey validates window=; empty defaults to overall.
func ParseHeadwayWindowKey(s string) (string, error) {
	if s == "" {
		return business.HeadwayWindowOverall, nil
	}
	if _, ok := ValidHeadwayWindowKeys[s]; !ok {
		return "", fmt.Errorf("invalid window %q", s)
	}
	return s, nil
}

// WindowsForServiceDate returns the window keys that apply on calendar date d (Chicago).
func WindowsForServiceDate(d time.Time) []HeadwayWindowSpec {
	loc := ChicagoLocation()
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
	dayEnd := day.AddDate(0, 0, 1)

	out := []HeadwayWindowSpec{{
		Key:   business.HeadwayWindowOverall,
		Start: day,
		End:   dayEnd,
	}}

	switch day.Weekday() {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		out = append(out,
			HeadwayWindowSpec{Key: business.HeadwayWindowWeekdayAM, Start: day.Add(6 * time.Hour), End: day.Add(9 * time.Hour)},
			HeadwayWindowSpec{Key: business.HeadwayWindowWeekdayMid, Start: day.Add(9 * time.Hour), End: day.Add(15 * time.Hour)},
			HeadwayWindowSpec{Key: business.HeadwayWindowWeekdayPM, Start: day.Add(15 * time.Hour), End: day.Add(19 * time.Hour)},
			HeadwayWindowSpec{Key: business.HeadwayWindowWeekdayEve, Start: day.Add(19 * time.Hour), End: dayEnd},
		)
	case time.Saturday:
		out = append(out, HeadwayWindowSpec{Key: business.HeadwayWindowSaturday, Start: day, End: dayEnd})
	case time.Sunday:
		out = append(out, HeadwayWindowSpec{Key: business.HeadwayWindowSunday, Start: day, End: dayEnd})
	}
	return out
}

// FilterHeadwaysByWindow keeps gaps whose later-arrival timestamp falls in [start, end).
func FilterHeadwaysByWindow(headways []business.Headway, spec HeadwayWindowSpec) []business.Headway {
	if len(headways) == 0 {
		return nil
	}
	out := make([]business.Headway, 0, len(headways))
	for _, h := range headways {
		ts := h.Timestamp.In(ChicagoLocation())
		if !ts.Before(spec.Start) && ts.Before(spec.End) {
			out = append(out, h)
		}
	}
	return out
}
