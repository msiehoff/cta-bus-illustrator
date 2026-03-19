import { useRef, useCallback, useEffect, useState } from 'react'
import Map, { Layer, Source, type MapRef } from 'react-map-gl/maplibre'
import type { LayerProps } from 'react-map-gl/maplibre'
import type { GetRoutesResponse } from '../types/api'
import 'maplibre-gl/dist/maplibre-gl.css'

// Positron is a minimal, light basemap — keeps focus on transit lines
const MAP_STYLE = 'https://tiles.openfreemap.org/styles/positron'

const CHICAGO_CENTER = { longitude: -87.6298, latitude: 41.8781 }
const MAX_ZOOM = 15

// CTA blue used as a consistent route color — thickness carries the data story.
const ROUTE_COLOR = '#009BDE'

const routeLineLayer: LayerProps = {
  id: 'bus-routes',
  type: 'line',
  paint: {
    'line-color': ROUTE_COLOR,
    // Log-spaced stops across Chicago's ridership range (~100–30k riders/day).
    // Coalesce to 0 so routes with no ridership data render at minimum width.
    'line-width': [
      'interpolate', ['linear'],
      ['coalesce', ['get', 'avgRides'], 0],
        0,  1,
     1000,  2,
     5000,  4,
    15000,  7,
    30000, 10,
    ],
    'line-opacity': 0.85,
  },
}

const RouteMap = () => {
  const mapRef = useRef<MapRef>(null)
  const [routes, setRoutes] = useState<GetRoutesResponse | null>(null)

  useEffect(() => {
    const fetchRoutes = async () => {
      const res = await fetch('/api/v1/routes')
      const data: GetRoutesResponse = await res.json()
      setRoutes(data)
    }
    fetchRoutes()
  }, [])

  const onMapLoad = useCallback(() => {
    mapRef.current?.resize()
  }, [])

  return (
    <Map
      ref={mapRef}
      initialViewState={{
        ...CHICAGO_CENTER,
        zoom: 11,
      }}
      style={{ width: '100%', height: '100%' }}
      mapStyle={MAP_STYLE}
      maxZoom={MAX_ZOOM}
      onLoad={onMapLoad}
    >
      {routes && (
        <Source id="bus-routes" type="geojson" data={routes}>
          <Layer {...routeLineLayer} />
        </Source>
      )}
    </Map>
  )
}

export default RouteMap
