export type RidershipType = 'weekday' | 'saturday' | 'sunday'

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
