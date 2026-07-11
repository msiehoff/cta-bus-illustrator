package business

import "strings"

// NormalizeDirection maps CTA direction strings (e.g. "EAST", "East Bound")
// onto the canonical forms used by getstops ("Eastbound", etc.).
func NormalizeDirection(rtdir string) string {
	s := strings.ToLower(strings.TrimSpace(rtdir))
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	switch {
	case strings.Contains(s, "north"):
		return "Northbound"
	case strings.Contains(s, "south"):
		return "Southbound"
	case strings.Contains(s, "east"):
		return "Eastbound"
	case strings.Contains(s, "west"):
		return "Westbound"
	default:
		return strings.TrimSpace(rtdir)
	}
}
