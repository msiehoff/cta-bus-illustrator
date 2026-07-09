import type { RouteComparison, RidershipDataPoint } from '../types/api'
import { EXPRESS_PAIRS, getPairedRouteId, isExpressRoute } from './expressPairs'

export const CORRIDOR_ROUTE_PREFIX = 'corridor-'

export const LOCAL_CORRIDOR_IDS = Object.keys(EXPRESS_PAIRS).filter(id => !isExpressRoute(id))

export const buildCorridorRouteId = (localId: string): string =>
  `${CORRIDOR_ROUTE_PREFIX}${localId}`

export const parseCorridorLocalId = (routeId: string): string | undefined =>
  routeId.startsWith(CORRIDOR_ROUTE_PREFIX)
    ? routeId.slice(CORRIDOR_ROUTE_PREFIX.length)
    : undefined

export const isCorridorRouteId = (routeId: string): boolean =>
  routeId.startsWith(CORRIDOR_ROUTE_PREFIX)

export const getCorridorDisplayName = (localRouteName: string): string =>
  `${localRouteName.replace(/ Express$/i, '')} Corridor`

export const getCorridorRoutePath = (localId: string): string =>
  `/routes/corridor/${localId}`

export const getCorridorBadgeLabel = (localId: string): string => {
  const expressId = getPairedRouteId(localId)
  return expressId ? `${localId}+${expressId}` : localId
}

const sumOptional = (a?: number, b?: number): number | undefined => {
  if (a == null && b == null) return undefined
  return (a ?? 0) + (b ?? 0)
}

const recoveryPct = (current: number, baseline: number): number | undefined => {
  if (baseline === 0) return undefined
  return (current / baseline) * 100
}

const pctChange = (current: number, baseline: number): number | undefined => {
  if (baseline === 0) return undefined
  return ((current - baseline) / baseline) * 100
}

export const mergeRouteComparisons = (
  local: RouteComparison,
  express: RouteComparison,
): RouteComparison => {
  const current = local.current + express.current
  const yearAgo = sumOptional(local.yearAgo, express.yearAgo)
  const fiveYearsAgo = sumOptional(local.fiveYearsAgo, express.fiveYearsAgo)
  const preCovid2019 = sumOptional(local.preCovid2019, express.preCovid2019)

  return {
    routeId: buildCorridorRouteId(local.routeId),
    routeName: getCorridorDisplayName(local.routeName),
    current,
    yearAgo,
    fiveYearsAgo,
    preCovid2019,
    recoveryPct: preCovid2019 != null ? recoveryPct(current, preCovid2019) : undefined,
    yearAgoPct: yearAgo != null ? pctChange(current, yearAgo) : undefined,
    fiveYearPct: fiveYearsAgo != null ? pctChange(current, fiveYearsAgo) : undefined,
  }
}

export const buildCorridorRows = (routes: RouteComparison[]): RouteComparison[] => {
  const byId = new Map(routes.map(r => [r.routeId, r]))
  const rows: RouteComparison[] = []

  for (const localId of LOCAL_CORRIDOR_IDS) {
    const expressId = getPairedRouteId(localId)
    if (!expressId) continue
    const local = byId.get(localId)
    const express = byId.get(expressId)
    if (local && express) {
      rows.push(mergeRouteComparisons(local, express))
    }
  }

  return rows
}

export const appendCorridorRows = (routes: RouteComparison[]): RouteComparison[] =>
  [...routes, ...buildCorridorRows(routes)]

export const mergeRidershipRecords = (
  localRecords: RidershipDataPoint[],
  expressRecords: RidershipDataPoint[],
): RidershipDataPoint[] => {
  const merged = new Map<string, RidershipDataPoint>()

  for (const rec of [...localRecords, ...expressRecords]) {
    const key = `${rec.month}|${rec.type}`
    const existing = merged.get(key)
    if (existing) {
      merged.set(key, { ...existing, avgRides: existing.avgRides + rec.avgRides })
    } else {
      merged.set(key, { ...rec })
    }
  }

  return [...merged.values()].sort(
    (a, b) => a.month.localeCompare(b.month) || a.type.localeCompare(b.type),
  )
}

export const matchesCorridorSearch = (
  route: RouteComparison,
  search: string,
): boolean => {
  if (!isCorridorRouteId(route.routeId)) return false
  const localId = parseCorridorLocalId(route.routeId)
  if (!localId) return false
  const expressId = getPairedRouteId(localId)
  const searchLower = search.toLowerCase()
  return (
    localId.toLowerCase().includes(searchLower) ||
    (expressId?.toLowerCase().includes(searchLower) ?? false) ||
    searchLower.includes('corridor')
  )
}
