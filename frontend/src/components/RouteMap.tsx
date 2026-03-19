import { useRef, useCallback, useEffect, useState } from 'react'
import Map, { Layer, Source, type MapRef, type MapLayerMouseEvent } from 'react-map-gl/maplibre'
import type { LayerProps } from 'react-map-gl/maplibre'
import type { GetRoutesResponse, RouteProperties } from '../types/api'
import RouteTooltip from './RouteTooltip'
import 'maplibre-gl/dist/maplibre-gl.css'

const MAP_STYLE = 'https://tiles.openfreemap.org/styles/positron'
const CHICAGO_CENTER = { longitude: -87.6298, latitude: 41.8781 }
const MAX_ZOOM = 15
const ROUTE_COLOR = '#009BDE'

// Maps each route id to its express/local counterpart.
// Both directions are listed so a lookup always works regardless of which is hovered.
const EXPRESS_PAIRS: Record<string, string> = {
  '4': 'X4', 'X4': '4',
  '9': 'X9', 'X9': '9',
  '49': 'X49', 'X49': '49',
}

const routeLineLayer: LayerProps = {
  id: 'bus-routes',
  type: 'line',
  paint: {
    'line-color': ROUTE_COLOR,
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

export interface SingleTooltipData {
  type: 'single'
  properties: RouteProperties
}

export interface CorridorTooltipData {
  type: 'corridor'
  local: RouteProperties
  express: RouteProperties
}

export type TooltipData = SingleTooltipData | CorridorTooltipData

interface HoveredRoute {
  data: TooltipData
  x: number
  y: number
}

const findRouteProperties = (routeId: string, routes: GetRoutesResponse): RouteProperties | null =>
  (routes.features.find(f => f.properties.routeId === routeId)?.properties) ?? null

const RouteMap = () => {
  const mapRef = useRef<MapRef>(null)
  const [routes, setRoutes] = useState<GetRoutesResponse | null>(null)
  const [hoveredRoute, setHoveredRoute] = useState<HoveredRoute | null>(null)

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

  const onMouseMove = useCallback((e: MapLayerMouseEvent) => {
    const canvas = mapRef.current?.getCanvas()
    const feature = e.features?.[0]

    if (!feature?.properties) {
      if (canvas) canvas.style.cursor = ''
      setHoveredRoute(null)
      return
    }

    if (canvas) canvas.style.cursor = 'pointer'
    const props = feature.properties as RouteProperties
    const { x, y } = e.point

    const pairedId = EXPRESS_PAIRS[props.routeId]
    if (pairedId && routes) {
      const pairedProps = findRouteProperties(pairedId, routes)
      if (pairedProps) {
        const isExpress = props.routeId.startsWith('X')
        setHoveredRoute({
          data: {
            type: 'corridor',
            local: isExpress ? pairedProps : props,
            express: isExpress ? props : pairedProps,
          },
          x,
          y,
        })
        return
      }
    }

    setHoveredRoute({ data: { type: 'single', properties: props }, x, y })
  }, [routes])

  const onMouseLeave = useCallback(() => {
    const canvas = mapRef.current?.getCanvas()
    if (canvas) canvas.style.cursor = ''
    setHoveredRoute(null)
  }, [])

  return (
    <div className="relative w-full h-full">
      <Map
        ref={mapRef}
        initialViewState={{ ...CHICAGO_CENTER, zoom: 11 }}
        style={{ width: '100%', height: '100%' }}
        mapStyle={MAP_STYLE}
        maxZoom={MAX_ZOOM}
        interactiveLayerIds={['bus-routes']}
        onLoad={onMapLoad}
        onMouseMove={onMouseMove}
        onMouseLeave={onMouseLeave}
      >
        {routes && (
          <Source id="bus-routes" type="geojson" data={routes}>
            <Layer {...routeLineLayer} />
          </Source>
        )}
      </Map>

      {hoveredRoute && <RouteTooltip data={hoveredRoute.data} x={hoveredRoute.x} y={hoveredRoute.y} />}
    </div>
  )
}

export default RouteMap
