import { useRef, useCallback } from 'react'
import Map, { Layer, Source, type MapRef } from 'react-map-gl/maplibre'
import type { LayerProps } from 'react-map-gl/maplibre'
import 'maplibre-gl/dist/maplibre-gl.css'
import { fakeRoutes } from '../data/routes'

// OpenFreeMap provides a free, open-source MapLibre-compatible tile style
const MAP_STYLE = 'https://tiles.openfreemap.org/styles/liberty'

const CHICAGO_CENTER = { longitude: -87.6298, latitude: 41.8781 }

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
      onLoad={onMapLoad}
    >
      <Source id="bus-routes" type="geojson" data={fakeRoutes}>
        <Layer {...routeLineLayer} />
      </Source>
    </Map>
  )
}

export default RouteMap
