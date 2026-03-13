package api

// fakeRoutes holds static GeoJSON until real data is served from the DB.
var fakeRoutes = map[string]any{
	"type": "FeatureCollection",
	"features": []map[string]any{
		{
			"type": "Feature",
			"properties": map[string]any{
				"routeId":   "66",
				"routeName": "Chicago",
				"color":     "#E63946",
				"ridership": 8500, // avg weekday riders
			},
			"geometry": map[string]any{
				"type": "LineString",
				"coordinates": [][]float64{
					{-87.8065, 41.8957}, // Austin & Chicago
					{-87.7754, 41.8957}, // Central & Chicago
					{-87.7523, 41.8957}, // Laramie & Chicago
					{-87.7341, 41.8957}, // Cicero & Chicago
					{-87.7137, 41.8957}, // Pulaski & Chicago
					{-87.6948, 41.8957}, // Kedzie & Chicago
					{-87.6726, 41.8957}, // Western & Chicago
					{-87.6558, 41.8957}, // Damen & Chicago
					{-87.6418, 41.8957}, // Ashland & Chicago
					{-87.6264, 41.8957}, // Halsted & Chicago
					{-87.6134, 41.8957}, // Morgan & Chicago
					{-87.6062, 41.8957}, // Peoria & Chicago
					{-87.5961, 41.8966}, // State & Chicago
					{-87.5854, 41.8978}, // Michigan & Chicago
					{-87.5769, 41.8983}, // McClurg & Chicago (Navy Pier turn)
				},
			},
		},
		{
			"type": "Feature",
			"properties": map[string]any{
				"routeId":   "66",
				"routeName": "Chicago",
				"color":     "#E63946",
				"ridership": 15000, // avg weekday riders
			},
			"geometry": map[string]any{
				"type": "LineString",
				"coordinates": [][]float64{

					{-87.8065, 41.95425369582809}, // Austin & Chicago
					{-87.7754, 41.95425369582809}, // Central & Chicago
					{-87.7523, 41.95425369582809}, // Laramie & Chicago
					{-87.7341, 41.95425369582809}, // Cicero & Chicago
					{-87.7137, 41.95425369582809}, // Pulaski & Chicago
					{-87.6948, 41.95425369582809}, // Kedzie & Chicago
					{-87.6726, 41.95425369582809}, // Western & Chicago
					{-87.6558, 41.95425369582809}, // Damen & Chicago
					{-87.6418, 41.95425369582809}, // Ashland & Chicago
					{-87.6264, 41.95425369582809}, // Halsted & Chicago
					{-87.6134, 41.95425369582809}, // Morgan & Chicago
					{-87.6062, 41.95425369582809}, // Peoria & Chicago
					{-87.5961, 41.95425369582809}, // State & Chicago
					{-87.5854, 41.95425369582809}, // Michigan & Chicago
					{-87.5769, 41.95425369582809}, // McClurg & Chicago (Navy Pier turn)
				},
			},
		},
	},
}
