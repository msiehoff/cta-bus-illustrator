export const formatHeadwayMinutes = (mins: number) => {
  if (!Number.isFinite(mins)) return '—'
  if (mins >= 10) return mins.toFixed(0)
  return mins.toFixed(1)
}

export const formatHeadwayCV = (cv: number) => {
  if (!Number.isFinite(cv)) return '—'
  return cv.toFixed(2)
}

/** User-facing name for coefficient of variation on public pages. */
export const HEADWAY_CONSISTENCY_LABEL = 'Consistency'

export const HEADWAY_CONSISTENCY_TOOLTIP =
  'How evenly buses arrive (also called CV). Lower is more clockwork; higher means a mix of bunches and long gaps.'

/** Plain-language band for a CV value. */
export const describeHeadwayConsistency = (
  cv: number,
): { label: string; trendUp?: boolean } => {
  if (!Number.isFinite(cv)) return { label: '—' }
  if (cv < 0.3) return { label: 'More even arrivals', trendUp: true }
  if (cv < 0.6) return { label: 'Typical variation' }
  return { label: 'Uneven — bunches & gaps', trendUp: false }
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

export interface HeadwayDensityResult {
  bins: HeadwayDensityPoint[]
  /** Routes with median above the chart axis max (outliers). */
  excludedCount: number
  axisMaxMinutes: number
}

const HEADWAY_DIST_ABS_CAP_MINUTES = 60

const percentileSorted = (sortedAsc: number[], p: number): number => {
  if (!sortedAsc.length) return 0
  const rank = Math.ceil((p / 100) * sortedAsc.length) - 1
  return sortedAsc[Math.min(sortedAsc.length - 1, Math.max(0, rank))]
}

/**
 * Build a density chart domain that ignores extreme outlier medians
 * (e.g. 400+ min) so typical routes stay readable.
 */
export const buildHeadwayDensity = (
  routes: { medianMinutes: number }[],
  binCount = 16,
): HeadwayDensityResult => {
  const values = routes
    .map(r => r.medianMinutes)
    .filter(v => Number.isFinite(v) && v > 0)
    .sort((a, b) => a - b)

  if (!values.length) {
    return { bins: [], excludedCount: 0, axisMaxMinutes: 0 }
  }

  const p95 = percentileSorted(values, 95)
  const axisMaxMinutes = Math.min(Math.max(p95, values[0]), HEADWAY_DIST_ABS_CAP_MINUTES)
  const clipped = values.filter(v => v <= axisMaxMinutes)
  const excludedCount = values.length - clipped.length

  if (!clipped.length) {
    return { bins: [], excludedCount, axisMaxMinutes }
  }

  if (clipped.length === 1) {
    return {
      bins: [{
        minutes: clipped[0],
        routeCount: 1,
        rangeLabel: `${formatHeadwayMinutes(clipped[0])} min`,
      }],
      excludedCount,
      axisMaxMinutes,
    }
  }

  const min = Math.min(...clipped)
  const max = Math.max(...clipped)
  const binWidth = (max - min) / binCount || 1

  const bins = Array.from({ length: binCount }, (_, i) => {
    const binMin = min + i * binWidth
    const binMax = min + (i + 1) * binWidth
    return { binMin, binMax, routeCount: 0 }
  })

  for (const value of clipped) {
    let index = Math.floor((value - min) / binWidth)
    if (index >= binCount) index = binCount - 1
    if (index < 0) index = 0
    bins[index].routeCount++
  }

  return {
    bins: bins.map(bin => ({
      minutes: (bin.binMin + bin.binMax) / 2,
      routeCount: bin.routeCount,
      rangeLabel: `${formatHeadwayMinutes(bin.binMin)}–${formatHeadwayMinutes(bin.binMax)} min`,
    })),
    excludedCount,
    axisMaxMinutes,
  }
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
