package business

type Route struct {
	ExternalID string
	Name       string

	Segments []RouteSegment
}

type RouteSegment struct {
	Sequence int
	Lat      float64
	Lng      float64
}
