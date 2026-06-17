import type { RidershipDataPoint, RidershipType } from '../types/api'

export const formatRides = (value: number): string => {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(1)}k`
  return value.toFixed(0)
}

export const formatMonth = (month: string): string => {
  const [year, mon] = month.split('-')
  return new Date(Number(year), Number(mon) - 1, 1)
    .toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
}

export const pctDiff = (current: number, baseline: number): number | null => {
  if (!baseline) return null
  return ((current - baseline) / baseline) * 100
}

export const recoveryPct = (current: number, baseline: number): number | null => {
  if (!baseline) return null
  return (current / baseline) * 100
}

export const formatPct = (pct: number): string =>
  `${pct >= 0 ? '+' : ''}${pct.toFixed(1)}%`

export const formatRecoveryPct = (pct: number): string =>
  `${pct.toFixed(0)}%`

export const comparisonMonth = (latestMonth: string, yearsBack: number): string => {
  const [year, mon] = latestMonth.split('-').map(Number)
  return `${year - yearsBack}-${String(mon).padStart(2, '0')}`
}

export const preCovidMonth = (latestMonth: string): string => {
  const [, mon] = latestMonth.split('-')
  return `2019-${mon}`
}

export const getRidership = (
  records: RidershipDataPoint[],
  month: string,
  type: RidershipType,
): number | null => {
  const row = records.find(r => r.month === month && r.type === type)
  return row ? row.avgRides : null
}

export const recoveryColorClass = (pct: number): string => {
  if (pct >= 100) return 'text-blue-400'
  if (pct >= 90) return 'text-green-400'
  if (pct >= 70) return 'text-amber-400'
  return 'text-red-400'
}

export const recoveryBarColorClass = (pct: number): string => {
  if (pct >= 100) return 'bg-blue-500'
  if (pct >= 90) return 'bg-green-500'
  if (pct >= 70) return 'bg-amber-500'
  return 'bg-red-500'
}

export const preCovidTrend = (
  records: RidershipDataPoint[],
  latestMonth: string,
  type: RidershipType,
): { trend?: string; trendUp?: boolean } => {
  const current = getRidership(records, latestMonth, type)
  const benchmark = preCovidMonth(latestMonth)
  const baseline = getRidership(records, benchmark, type)
  if (current == null || baseline == null) return {}
  const pct = recoveryPct(current, baseline)
  if (pct == null) return {}
  return {
    trend: `${formatRecoveryPct(pct)} of ${formatMonth(benchmark)}`,
    trendUp: pct >= 100 ? true : pct < 70 ? false : undefined,
  }
}

export interface RecoveryRow {
  label: string
  type: RidershipType
  current: number
  preCovid2019: number | null
  yearAgo: number | null
  fiveYearsAgo: number | null
  recoveryPct: number | null
  yearAgoPct: number | null
  fiveYearPct: number | null
}

export interface RouteRecoverySnapshot {
  currentMonth: string
  benchmarkMonth: string
  yearAgoMonth: string
  fiveYearsAgoMonth: string
  rows: RecoveryRow[]
}

const DAY_TYPES: { label: string; type: RidershipType }[] = [
  { label: 'Weekday', type: 'weekday' },
  { label: 'Saturday', type: 'saturday' },
  { label: 'Sunday', type: 'sunday' },
]

export const buildRouteRecovery = (
  records: RidershipDataPoint[],
  latestMonth: string,
): RouteRecoverySnapshot => {
  const benchmarkMonth = preCovidMonth(latestMonth)
  const yearAgoMonth = comparisonMonth(latestMonth, 1)
  const fiveYearsAgoMonth = comparisonMonth(latestMonth, 5)

  const rows = DAY_TYPES.map(({ label, type }) => {
    const current = getRidership(records, latestMonth, type) ?? 0
    const preCovid2019 = getRidership(records, benchmarkMonth, type)
    const yearAgo = getRidership(records, yearAgoMonth, type)
    const fiveYearsAgo = getRidership(records, fiveYearsAgoMonth, type)

    return {
      label,
      type,
      current,
      preCovid2019,
      yearAgo,
      fiveYearsAgo,
      recoveryPct: preCovid2019 != null ? recoveryPct(current, preCovid2019) : null,
      yearAgoPct: yearAgo != null ? pctDiff(current, yearAgo) : null,
      fiveYearPct: fiveYearsAgo != null ? pctDiff(current, fiveYearsAgo) : null,
    }
  })

  return { currentMonth: latestMonth, benchmarkMonth, yearAgoMonth, fiveYearsAgoMonth, rows }
}

export interface SeasonalityPoint {
  monthNum: number
  label: string
  avgRides: number
  sampleYears: number
}

const MONTH_LABELS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

export const buildSeasonality = (
  records: RidershipDataPoint[],
  type: RidershipType,
): SeasonalityPoint[] => {
  const buckets = new Map<number, number[]>()

  for (const record of records) {
    if (record.type !== type) continue
    const monthNum = Number(record.month.split('-')[1])
    if (!buckets.has(monthNum)) buckets.set(monthNum, [])
    buckets.get(monthNum)!.push(record.avgRides)
  }

  return MONTH_LABELS.map((label, i) => {
    const monthNum = i + 1
    const values = buckets.get(monthNum) ?? []
    const avgRides = values.length
      ? values.reduce((sum, v) => sum + v, 0) / values.length
      : 0
    return { monthNum, label, avgRides, sampleYears: values.length }
  }).filter(point => point.sampleYears > 0)
}

export interface WeekendSharePoint {
  month: string
  sharePct: number
}

export const buildWeekendShare = (records: RidershipDataPoint[]): WeekendSharePoint[] => {
  const byMonth = new Map<string, Partial<Record<RidershipType, number>>>()

  for (const record of records) {
    if (!byMonth.has(record.month)) byMonth.set(record.month, {})
    byMonth.get(record.month)![record.type] = record.avgRides
  }

  return Array.from(byMonth.entries())
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([month, types]) => {
      const weekday = types.weekday ?? 0
      const weekend = (types.saturday ?? 0) + (types.sunday ?? 0)
      const sharePct = weekday > 0 ? (weekend / weekday) * 100 : 0
      return { month, sharePct }
    })
    .filter(point => point.sharePct > 0)
}

export interface RouteRankInfo {
  rank: number
  total: number
  networkSharePct: number
}

export const getRouteRank = (
  routeId: string,
  routes: { routeId: string; current: number }[],
  systemTotal: number,
): RouteRankInfo | null => {
  const sorted = [...routes].sort((a, b) => b.current - a.current)
  const index = sorted.findIndex(r => r.routeId === routeId)
  if (index < 0 || !systemTotal) return null
  return {
    rank: index + 1,
    total: sorted.length,
    networkSharePct: (sorted[index].current / systemTotal) * 100,
  }
}

export interface PeerRoute {
  routeId: string
  routeName: string
  current: number
  recoveryPct?: number
}

export const getPeerRoutes = (
  routeId: string,
  routes: { routeId: string; routeName: string; current: number; recoveryPct?: number }[],
  limit = 4,
): PeerRoute[] => {
  const target = routes.find(r => r.routeId === routeId)
  if (!target) return []

  const band = target.current * 0.2
  return routes
    .filter(r => r.routeId !== routeId && Math.abs(r.current - target.current) <= band)
    .sort((a, b) => Math.abs(a.current - target.current) - Math.abs(b.current - target.current))
    .slice(0, limit)
}

export interface RecoveryBucket {
  label: string
  count: number
}

export const buildRecoveryDistribution = (
  routes: { recoveryPct?: number }[],
): RecoveryBucket[] => {
  const buckets = [
    { label: '<70%', min: 0, max: 70, count: 0 },
    { label: '70–89%', min: 70, max: 90, count: 0 },
    { label: '90–99%', min: 90, max: 100, count: 0 },
    { label: '100%+', min: 100, max: Infinity, count: 0 },
  ]

  for (const route of routes) {
    if (route.recoveryPct == null) continue
    const bucket = buckets.find(b => route.recoveryPct! >= b.min && route.recoveryPct! < b.max)
      ?? buckets[buckets.length - 1]
    bucket.count++
  }

  return buckets.map(({ label, count }) => ({ label, count }))
}

export interface ScatterPoint {
  routeId: string
  routeName: string
  current: number
  recoveryPct: number
}

export const buildRecoveryScatter = (
  routes: { routeId: string; routeName: string; current: number; recoveryPct?: number }[],
): ScatterPoint[] =>
  routes
    .filter(r => r.recoveryPct != null && r.current > 0)
    .map(r => ({
      routeId: r.routeId,
      routeName: r.routeName,
      current: r.current,
      recoveryPct: r.recoveryPct!,
    }))

export const getSnapshotForMonth = (
  records: RidershipDataPoint[],
  month: string,
): { month: string; weekday: number; saturday: number; sunday: number } | null => {
  const rows = records.filter(r => r.month === month)
  if (!rows.length) return null
  const get = (type: RidershipType) => rows.find(r => r.type === type)?.avgRides ?? 0
  return {
    month,
    weekday: get('weekday'),
    saturday: get('saturday'),
    sunday: get('sunday'),
  }
}
