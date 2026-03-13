package business

type Route struct {
	ExternalID string
	Name       string

	Color     string
	Ridership int
	Segments  []RouteSegment
}

type RouteSegment struct {
	Lat float64
	Lng float64
}
