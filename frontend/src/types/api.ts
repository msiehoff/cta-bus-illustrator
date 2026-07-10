export type RidershipType = 'weekday' | 'saturday' | 'sunday'

// --- Map / routes endpoint ---

export interface RouteProperties {
  routeId: string
  routeName: string
  avgRides?: number
}

export interface RouteGeometry {
  type: 'LineString'
  coordinates: [number, number][]
}

export interface RouteFeature {
  type: 'Feature'
  properties: RouteProperties
  geometry: RouteGeometry
}

export interface GetRoutesResponse {
  type: 'FeatureCollection'
  features: RouteFeature[]
}

// --- Ridership time-series endpoints ---

export interface RidershipDataPoint {
  month: string       // "2024-01"
  type: RidershipType
  avgRides: number
}

export interface GetRidershipResponse {
  records: RidershipDataPoint[]
}

// --- Routes comparison endpoint ---

export interface RouteComparison {
  routeId: string
  routeName: string
  current: number
  yearAgo?: number
  fiveYearsAgo?: number
  preCovid2019?: number
  recoveryPct?: number
  yearAgoPct?: number
  fiveYearPct?: number
}

export interface GetRoutesComparisonResponse {
  currentMonth: string
  benchmarkMonth: string
  yearAgoMonth: string
  fiveYearsAgoMonth: string
  systemCurrent: number
  systemPreCovid?: number
  systemRecovery?: number
  routes: RouteComparison[]
}

// --- Admin endpoints ---

export interface AdminSessionResponse {
  authenticated: boolean
  username?: string
}

export interface PipelineStatusResponse {
  enabled: boolean
  running: boolean
  routes: string[]
  routeCount: number
  pollInterval: string
  lastPollAt?: string
  lastPingCount: number
  lastError?: string
  startedAt?: string
  arrivalCount: number
}

export interface ArrivalRecord {
  stopId: string
  routeId: string
  direction: string
  vehicleId: string
  timestamp: string
}

export interface ListArrivalsResponse {
  arrivals: ArrivalRecord[]
  total: number
  limit: number
  offset: number
}
