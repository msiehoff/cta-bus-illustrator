export interface RouteProperties {
  routeId: string
  routeName: string
  color: string
  ridership: number
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
