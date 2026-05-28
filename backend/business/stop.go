package business

type Stop struct {
	StopID    string
	RouteID   string
	Direction string
	Name      string
	Lat       float64
	Lon       float64
	Sequence  int
}
