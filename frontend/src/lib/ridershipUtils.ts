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
