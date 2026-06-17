export const EXPRESS_PAIRS: Record<string, string> = {
  '4': 'X4', 'X4': '4',
  '9': 'X9', 'X9': '9',
  '49': 'X49', 'X49': '49',
}

export const getPairedRouteId = (routeId: string): string | undefined =>
  EXPRESS_PAIRS[routeId]

export const isExpressRoute = (routeId: string): boolean =>
  routeId.startsWith('X')
