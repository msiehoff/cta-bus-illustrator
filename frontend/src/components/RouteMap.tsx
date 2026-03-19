import { useRef, useCallback, useEffect, useState, useMemo } from 'react'
import Map, { Layer, Source, type MapRef, type MapLayerMouseEvent } from 'react-map-gl/maplibre'
import type { LayerProps } from 'react-map-gl/maplibre'
import type { GetRoutesResponse, RouteProperties, RidershipType } from '../types/api'
import { ridershipColorExpression } from '../lib/ridershipColors'
import RouteTooltip from './RouteTooltip'
import FilterBar from './FilterBar'
import 'maplibre-gl/dist/maplibre-gl.css'

const MAP_STYLE = 'https://tiles.openfreemap.org/styles/positron'
const CHICAGO_CENTER = { longitude: -87.6298, latitude: 41.8781 }
const MAX_ZOOM = 15

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
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    'line-color': ridershipColorExpression as any,
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

export interface RankedEntry {
  rank: number
  data: TooltipData
  totalRides: number
}

interface HoveredRoute {
  data: TooltipData
  x: number
  y: number
}

const findRouteProperties = (routeId: string, routes: GetRoutesResponse): RouteProperties | null =>
  (routes.features.find(f => f.properties.routeId === routeId)?.properties) ?? null

const computeRankedEntries = (routes: GetRoutesResponse): RankedEntry[] => {
  const seen = new Set<string>()
  const entries: Array<{ data: TooltipData; totalRides: number }> = []

  for (const feature of routes.features) {
    const props = feature.properties
    if (seen.has(props.routeId)) continue
    seen.add(props.routeId)

    const pairedId = EXPRESS_PAIRS[props.routeId]
    if (pairedId) {
      const pairedProps = findRouteProperties(pairedId, routes)
      if (pairedProps) {
        seen.add(pairedId)
        const isExpress = props.routeId.startsWith('X')
        const local = isExpress ? pairedProps : props
        const express = isExpress ? props : pairedProps
        entries.push({
          data: { type: 'corridor', local, express },
          totalRides: (local.avgRides ?? 0) + (express.avgRides ?? 0),
        })
        continue
      }
    }

    entries.push({
      data: { type: 'single', properties: props },
      totalRides: props.avgRides ?? 0,
    })
  }

  return entries
    .filter(e => e.totalRides > 0)
    .sort((a, b) => b.totalRides - a.totalRides)
    .slice(0, 10)
    .map((entry, i) => ({ ...entry, rank: i + 1 }))
}

const RouteMap = () => {
  const mapRef = useRef<MapRef>(null)
  const [routes, setRoutes] = useState<GetRoutesResponse | null>(null)
  const [hoveredRoute, setHoveredRoute] = useState<HoveredRoute | null>(null)

  const [availableMonths, setAvailableMonths] = useState<string[]>([])
  const [selectedMonth, setSelectedMonth] = useState<string | null>(null)
  const [ridershipType, setRidershipType] = useState<RidershipType>('weekday')
  const [monthsLoaded, setMonthsLoaded] = useState(false)

  // Highlight state: ids = route IDs to illuminate, key = timestamp to re-trigger same route.
  const [highlight, setHighlight] = useState<{ ids: string[]; key: number } | null>(null)
  const [highlightVisible, setHighlightVisible] = useState(true)

  // Top 10 ranked corridors/routes derived from current ridership data.
  const rankedEntries = useMemo(() => routes ? computeRankedEntries(routes) : [], [routes])

  // Quick lookup: routeId → rank (both IDs of a corridor pair map to the same rank).
  const rankByRouteId = useMemo(() => {
    const map: Record<string, number> = {}
    for (const entry of rankedEntries) {
      if (entry.data.type === 'single') {
        map[entry.data.properties.routeId] = entry.rank
      } else {
        map[entry.data.local.routeId] = entry.rank
        map[entry.data.express.routeId] = entry.rank
      }
    }
    return map
  }, [rankedEntries])

  // Blink the highlight layer 4 times over ~2 seconds, then clear.
  useEffect(() => {
    if (!highlight) return
    setHighlightVisible(true)
    let count = 0
    const interval = setInterval(() => {
      count++
      setHighlightVisible(v => !v)
      if (count >= 8) {
        clearInterval(interval)
        setHighlight(null)
        setHighlightVisible(true)
      }
    }, 250)
    return () => clearInterval(interval)
  }, [highlight?.key]) // key changes on every click, so same route can re-trigger

  const onRouteListClick = useCallback((entry: RankedEntry) => {
    const ids = entry.data.type === 'single'
      ? [entry.data.properties.routeId]
      : [entry.data.local.routeId, entry.data.express.routeId]
    setHighlight({ ids, key: Date.now() })
  }, [])

  // Fetch available ridership months once on mount.
  useEffect(() => {
    const fetchMonths = async () => {
      const res = await fetch('/api/v1/ridership/months')
      const data = await res.json()
      const months: string[] = data.months ?? []
      setAvailableMonths(months)
      setSelectedMonth(months[0] ?? null)
      setMonthsLoaded(true)
    }
    fetchMonths()
  }, [])

  // Re-fetch routes whenever the selected month or ridership type changes.
  useEffect(() => {
    if (!monthsLoaded) return
    const fetchRoutes = async () => {
      const params = new URLSearchParams({ type: ridershipType })
      if (selectedMonth) params.set('month', selectedMonth)
      const res = await fetch(`/api/v1/routes?${params}`)
      const data: GetRoutesResponse = await res.json()
      setRoutes(data)
    }
    fetchRoutes()
  }, [monthsLoaded, selectedMonth, ridershipType])

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

  const hoveredRank = hoveredRoute
    ? hoveredRoute.data.type === 'single'
      ? rankByRouteId[hoveredRoute.data.properties.routeId]
      : rankByRouteId[hoveredRoute.data.local.routeId]
    : undefined

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
            {highlight && (
              <Layer
                id="bus-routes-highlight"
                type="line"
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                filter={['match', ['get', 'routeId'], highlight.ids, true, false] as any}
                paint={{
                  'line-color': '#ffffff',
                  'line-width': 14,
                  'line-opacity': highlightVisible ? 0.8 : 0,
                  // Disable MapLibre's default 300 ms transition so blinks are sharp.
                  // eslint-disable-next-line @typescript-eslint/no-explicit-any
                  'line-opacity-transition': { duration: 0, delay: 0 } as any,
                  'line-blur': 3,
                }}
              />
            )}
          </Source>
        )}
      </Map>

      <FilterBar
        availableMonths={availableMonths}
        selectedMonth={selectedMonth}
        ridershipType={ridershipType}
        rankedEntries={rankedEntries}
        onMonthChange={setSelectedMonth}
        onTypeChange={setRidershipType}
        onRouteClick={onRouteListClick}
      />

      {hoveredRoute && (
        <RouteTooltip
          data={hoveredRoute.data}
          x={hoveredRoute.x}
          y={hoveredRoute.y}
          rank={hoveredRank}
        />
      )}
    </div>
  )
}

export default RouteMap
