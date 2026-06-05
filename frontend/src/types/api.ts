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
