import type { Feature, LineString, FeatureCollection } from 'geojson'

export interface BusRouteProperties {
  routeId: string
  routeName: string
  color: string
}

// Fake route data for CTA Route 66 (Chicago Ave) — west to east across Chicago
const route66: Feature<LineString, BusRouteProperties> = {
  type: 'Feature',
  properties: {
    routeId: '66',
    routeName: 'Chicago',
    color: '#E63946',
  },
  geometry: {
    type: 'LineString',
    coordinates: [
      [-87.8065, 41.8957], // Austin & Chicago
      [-87.7754, 41.8957], // Central & Chicago
      [-87.7523, 41.8957], // Laramie & Chicago
      [-87.7341, 41.8957], // Cicero & Chicago
      [-87.7137, 41.8957], // Pulaski & Chicago
      [-87.6948, 41.8957], // Kedzie & Chicago
      [-87.6726, 41.8957], // Western & Chicago
      [-87.6558, 41.8957], // Damen & Chicago
      [-87.6418, 41.8957], // Ashland & Chicago
      [-87.6264, 41.8957], // Halsted & Chicago
      [-87.6134, 41.8957], // Morgan & Chicago
      [-87.6062, 41.8957], // Peoria & Chicago
      [-87.5961, 41.8966], // State & Chicago
      [-87.5854, 41.8978], // Michigan & Chicago
      [-87.5769, 41.8983], // McClurg & Chicago (Navy Pier turn)
    ],
  },
}

export const fakeRoutes: FeatureCollection<LineString, BusRouteProperties> = {
  type: 'FeatureCollection',
  features: [route66],
}
