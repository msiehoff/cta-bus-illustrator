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
  stopName?: string
  routeId: string
  routeName?: string
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

export interface HeadwayJobRun {
  id: number
  serviceDate: string
  status: string
  triggeredBy: string
  startedAt: string
  finishedAt?: string
  arrivalsProcessed: number
  headwaysWritten: number
  summariesWritten?: number
  errorMessage?: string
}

export interface ListHeadwayJobRunsResponse {
  runs: HeadwayJobRun[]
  limit: number
  offset: number
}

export interface HeadwayRecord {
  stopId: string
  stopName?: string
  routeId: string
  routeName?: string
  direction: string
  timestamp: string
  headwayMinutes: number
  fromVehicleId?: string
  toVehicleId?: string
}

export interface ListHeadwaysResponse {
  headways: HeadwayRecord[]
  total: number
  limit: number
  offset: number
}

export interface HeadwaySummaryStats {
  count: number
  meanMinutes: number
  medianMinutes: number
  stdDevMinutes: number
  cv: number
  avgWaitMinutes: number
}

export interface HeadwayStopSummary extends HeadwaySummaryStats {
  stopId: string
  stopName?: string
  routeId: string
  routeName?: string
  direction: string
}

export interface HeadwaySummaryResponse {
  pooled: HeadwaySummaryStats
  equalStopWeight: HeadwaySummaryStats
  byStop: HeadwayStopSummary[]
  source?: 'stored' | 'computed' | string
}

export interface HeadwaySummaryRow {
  serviceDate: string
  grain: string
  method: string
  stopId?: string
  stopName?: string
  routeId?: string
  routeName?: string
  direction?: string
  count: number
  meanMinutes: number
  medianMinutes: number
  stdDevMinutes: number
  cv: number
  avgWaitMinutes: number
}

export interface ListHeadwaySummariesResponse {
  summaries: HeadwaySummaryRow[]
  total: number
  limit: number
  offset: number
}

// --- Public headway endpoints ---

export interface HeadwayPeriodStats {
  count: number
  meanMinutes: number
  medianMinutes: number
  stdDevMinutes: number
  cv: number
  avgWaitMinutes: number
  daysWithData: number
  periodStart?: string
  periodEnd?: string
}

export interface HeadwayRoutePeriod extends HeadwayPeriodStats {
  routeId: string
  routeName?: string
}

export interface HeadwayDayPoint {
  serviceDate: string
  medianMinutes: number
  avgWaitMinutes: number
  cv: number
  count: number
}

export interface HeadwayRoutesListResponse {
  period: HeadwayPeriodStats
  routes: HeadwayRoutePeriod[]
  method: string
  grain: string
  days: number
}

export interface HeadwayRouteDetailResponse {
  route: HeadwayRoutePeriod
  series: HeadwayDayPoint[]
  method: string
  grain: string
  days: number
}

export interface HeadwaySystemResponse {
  period: HeadwayPeriodStats
  series: HeadwayDayPoint[]
  longestWaits: HeadwayRoutePeriod[]
  method: string
  grain: string
  days: number
}

