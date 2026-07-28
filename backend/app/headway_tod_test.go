package app_test

import (
	"testing"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

func TestWindowsForServiceDateWeekday(t *testing.T) {
	// 2026-07-10 is a Friday
	d, err := app.ParseServiceDate("2026-07-10")
	if err != nil {
		t.Fatal(err)
	}
	specs := app.WindowsForServiceDate(d)
	keys := make([]string, len(specs))
	for i, s := range specs {
		keys[i] = s.Key
	}
	want := []string{
		business.HeadwayWindowOverall,
		business.HeadwayWindowWeekdayAM,
		business.HeadwayWindowWeekdayMid,
		business.HeadwayWindowWeekdayPM,
		business.HeadwayWindowWeekdayEve,
	}
	if len(keys) != len(want) {
		t.Fatalf("keys = %v want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("keys[%d] = %q want %q", i, keys[i], want[i])
		}
	}
}

func TestWindowsForServiceDateSaturday(t *testing.T) {
	d, err := app.ParseServiceDate("2026-07-11") // Saturday
	if err != nil {
		t.Fatal(err)
	}
	specs := app.WindowsForServiceDate(d)
	if len(specs) != 2 {
		t.Fatalf("len = %d want 2", len(specs))
	}
	if specs[0].Key != business.HeadwayWindowOverall || specs[1].Key != business.HeadwayWindowSaturday {
		t.Fatalf("keys = %s,%s", specs[0].Key, specs[1].Key)
	}
}

func TestFilterHeadwaysByWindowAMOnly(t *testing.T) {
	loc := app.ChicagoLocation()
	day := time.Date(2026, 7, 10, 0, 0, 0, 0, loc)
	headways := []business.Headway{
		{StopID: "a", RouteID: "8", Direction: "NB", Timestamp: day.Add(7*time.Hour + 30*time.Minute), HeadwayMinutes: 10},
		{StopID: "a", RouteID: "8", Direction: "NB", Timestamp: day.Add(12 * time.Hour), HeadwayMinutes: 12},
		{StopID: "a", RouteID: "8", Direction: "NB", Timestamp: day.Add(16 * time.Hour), HeadwayMinutes: 14},
	}

	var amSpec app.HeadwayWindowSpec
	for _, s := range app.WindowsForServiceDate(day) {
		if s.Key == business.HeadwayWindowWeekdayAM {
			amSpec = s
			break
		}
	}
	got := app.FilterHeadwaysByWindow(headways, amSpec)
	if len(got) != 1 {
		t.Fatalf("am gaps = %d want 1", len(got))
	}
	if got[0].HeadwayMinutes != 10 {
		t.Fatalf("minutes = %.0f", got[0].HeadwayMinutes)
	}

	overall := app.WindowsForServiceDate(day)[0]
	all := app.FilterHeadwaysByWindow(headways, overall)
	if len(all) != 3 {
		t.Fatalf("overall gaps = %d want 3", len(all))
	}
}

func TestParseHeadwayWindowKey(t *testing.T) {
	k, err := app.ParseHeadwayWindowKey("")
	if err != nil || k != business.HeadwayWindowOverall {
		t.Fatalf("empty: %q %v", k, err)
	}
	k, err = app.ParseHeadwayWindowKey("weekday_am")
	if err != nil || k != business.HeadwayWindowWeekdayAM {
		t.Fatalf("am: %q %v", k, err)
	}
	if _, err := app.ParseHeadwayWindowKey("nope"); err == nil {
		t.Fatal("expected error")
	}
}
