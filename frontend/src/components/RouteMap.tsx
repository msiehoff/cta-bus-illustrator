import { useRef, useCallback, useEffect, useState } from 'react'
import Map, { Layer, Source, type MapRef } from 'react-map-gl/maplibre'
import type { LayerProps } from 'react-map-gl/maplibre'
import type { FeatureCollection } from 'geojson'
import 'maplibre-gl/dist/maplibre-gl.css'

// Positron is a minimal, light basemap — keeps focus on transit lines
const MAP_STYLE = 'https://tiles.openfreemap.org/styles/positron'

const CHICAGO_CENTER = { longitude: -87.6298, latitude: 41.8781 }
const MAX_ZOOM = 15

const routeLineLayer: LayerProps = {
  id: 'bus-routes',
  type: 'line',
  paint: {
    'line-color': ['get', 'color'],
    'line-width': 4,
    'line-opacity': 0.9,
  },
}

const RouteMap = () => {
  const mapRef = useRef<MapRef>(null)
  const [routes, setRoutes] = useState<FeatureCollection | null>(null)

  useEffect(() => {
    const fetchRoutes = async () => {
      const res = await fetch('/api/v1/routes')
      const data: FeatureCollection = await res.json()
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
