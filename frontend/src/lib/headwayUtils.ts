export const formatHeadwayMinutes = (mins: number) => {
  if (!Number.isFinite(mins)) return '—'
  if (mins >= 10) return mins.toFixed(0)
  return mins.toFixed(1)
}

export const formatHeadwayCV = (cv: number) => {
  if (!Number.isFinite(cv)) return '—'
  return cv.toFixed(2)
}

export const formatHeadwayPeriod = (start?: string, end?: string, daysWithData?: number) => {
  if (!start || !end) {
    return daysWithData ? `${daysWithData} day${daysWithData === 1 ? '' : 's'} of data` : 'No data yet'
  }
  if (start === end) return start
  return `${start} → ${end}`
}

export interface HeadwayRankInfo {
  rank: number
  total: number
}

/** Rank 1 = shortest median headway (better service). */
export const getHeadwayRank = (
  routeId: string,
  routes: { routeId: string; medianMinutes: number }[],
): HeadwayRankInfo | null => {
  const eligible = routes.filter(r => r.medianMinutes > 0)
  const sorted = [...eligible].sort((a, b) => a.medianMinutes - b.medianMinutes)
  const index = sorted.findIndex(r => r.routeId === routeId)
  if (index < 0) return null
  return { rank: index + 1, total: sorted.length }
}

/** Percent vs network: negative = shorter than network (better). */
export const headwayVsNetworkPct = (routeMinutes: number, networkMinutes: number): number | null => {
  if (!Number.isFinite(routeMinutes) || !Number.isFinite(networkMinutes) || networkMinutes <= 0) {
    return null
  }
  return ((routeMinutes - networkMinutes) / networkMinutes) * 100
}

export interface HeadwayDensityPoint {
  minutes: number
  routeCount: number
  rangeLabel: string
}

export const buildHeadwayDensity = (
  routes: { medianMinutes: number }[],
  binCount = 16,
): HeadwayDensityPoint[] => {
  const values = routes.map(r => r.medianMinutes).filter(v => v > 0)
  if (!values.length) return []
  if (values.length === 1) {
    return [{
      minutes: values[0],
      routeCount: 1,
      rangeLabel: `${formatHeadwayMinutes(values[0])} min`,
    }]
  }

  const min = Math.min(...values)
  const max = Math.max(...values)
  const binWidth = (max - min) / binCount || 1

  const bins = Array.from({ length: binCount }, (_, i) => {
    const binMin = min + i * binWidth
    const binMax = min + (i + 1) * binWidth
    return { binMin, binMax, routeCount: 0 }
  })

  for (const value of values) {
    let index = Math.floor((value - min) / binWidth)
    if (index >= binCount) index = binCount - 1
    if (index < 0) index = 0
    bins[index].routeCount++
  }

  return bins.map(bin => ({
    minutes: (bin.binMin + bin.binMax) / 2,
    routeCount: bin.routeCount,
    rangeLabel: `${formatHeadwayMinutes(bin.binMin)}–${formatHeadwayMinutes(bin.binMax)} min`,
  }))
}

/** Share of routes with a longer median than `value` (higher = this route is relatively frequent). */
export const getShorterThanPercentile = (
  routes: { medianMinutes: number }[],
  value: number,
): number => {
  if (!routes.length) return 0
  const longer = routes.filter(r => r.medianMinutes > value).length
  return Math.round((longer / routes.length) * 100)
}
