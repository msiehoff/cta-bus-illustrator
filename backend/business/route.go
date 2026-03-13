package business

type Route struct {
	ExternalID string
	Name       string

	Color       string
	Ridership   int
	Coordinates [][2]float64
}
