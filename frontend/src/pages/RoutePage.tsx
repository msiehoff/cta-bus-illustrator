import { useMemo, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useRouteRidership } from '../hooks/useRouteRidership'
import RidershipChart, { type WindowKey, cutoffMonth } from '../components/RidershipChart'
import StatCard from '../components/StatCard'
import type { RidershipType } from '../types/api'

function formatRides(value: number): string {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(1)}k`
  return value.toFixed(0)
}

function formatMonth(month: string): string {
  const [year, mon] = month.split('-')
  return new Date(Number(year), Number(mon) - 1, 1)
    .toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
}

function pctDiff(current: number, baseline: number): number | null {
  if (!baseline) return null
  return ((current - baseline) / baseline) * 100
}

function formatPct(pct: number): string {
  return `${pct >= 0 ? '+' : ''}${pct.toFixed(1)}%`
}

export default function RoutePage() {
  const { externalId = '' } = useParams<{ externalId: string }>()
  const navigate = useNavigate()
  const { records, loading, error } = useRouteRidership(externalId)
  const [window, setWindow] = useState<WindowKey>('5y')

  const routeName = (history.state as { routeName?: string } | null)?.routeName ?? externalId

  // Latest month in the dataset — used as the anchor for window cutoffs
  const latestMonth = useMemo(() => {
    if (!records.length) return null
    return records[records.length - 1].month
  }, [records])

  // Latest month values (the "current" numbers shown on stat cards)
  const latest = useMemo(() => {
    if (!latestMonth) return null
    const rows = records.filter(r => r.month === latestMonth)
    const get = (type: RidershipType) => rows.find(r => r.type === type)?.avgRides ?? 0
    return {
      month: latestMonth,
      weekday:  get('weekday'),
      saturday: get('saturday'),
      sunday:   get('sunday'),
    }
  }, [records, latestMonth])

  // First month within the current window — baseline for % comparison.
  // Cutoff is anchored to the latest data month, not today.
  const windowStart = useMemo(() => {
    if (!records.length || !latestMonth) return null
    const cutoff = cutoffMonth(window, latestMonth)
    const windowRecords = cutoff
      ? records.filter(r => r.month >= cutoff)
      : records
    if (!windowRecords.length) return null
    const firstMonth = windowRecords[0].month
    const rows = windowRecords.filter(r => r.month === firstMonth)
    const get = (type: RidershipType) => rows.find(r => r.type === type)?.avgRides ?? 0
    return {
      month: firstMonth,
      weekday:  get('weekday'),
      saturday: get('saturday'),
      sunday:   get('sunday'),
    }
  }, [records, window, latestMonth])

  const sinceLabel = windowStart ? `since ${formatMonth(windowStart.month)}` : null

  function statProps(type: 'weekday' | 'saturday' | 'sunday') {
    const current = latest?.[type] ?? 0
    const baseline = windowStart?.[type] ?? 0
    const pct = pctDiff(current, baseline)
    return {
      value: formatRides(current),
      trend: pct !== null && sinceLabel ? `${formatPct(pct)} ${sinceLabel}` : undefined,
      trendUp: pct !== null ? pct >= 0 : undefined,
    }
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      {/* Back */}
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-1.5 text-xs text-gray-500 hover:text-gray-300 mb-4 transition-colors"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <line x1="19" y1="12" x2="5" y2="12" />
          <polyline points="12 19 5 12 12 5" />
        </svg>
        System Overview
      </button>

      {/* Route header */}
      <div className="flex items-center gap-4 mb-5">
        <div className="flex items-center justify-center w-12 h-12 rounded-full bg-blue-900 text-blue-200 text-lg font-bold shrink-0">
          {externalId}
        </div>
        <div>
          <h1 className="text-xl font-semibold text-white">{routeName}</h1>
          {latest && (
            <p className="text-sm text-gray-400 mt-0.5">
              {formatRides(latest.weekday)} avg weekday riders · {formatMonth(latest.month)}
            </p>
          )}
        </div>
      </div>

      {/* Stat cards */}
      {latest && (
        <div className="grid grid-cols-3 gap-3 mb-5">
          <StatCard
            label={`Weekday riders · ${formatMonth(latest.month)}`}
            {...statProps('weekday')}
          />
          <StatCard
            label={`Saturday riders · ${formatMonth(latest.month)}`}
            {...statProps('saturday')}
          />
          <StatCard
            label={`Sunday riders · ${formatMonth(latest.month)}`}
            {...statProps('sunday')}
          />
        </div>
      )}

      {/* Chart */}
      <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-4">Ridership over time</h2>
        {loading && <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>}
        {error && <div className="text-red-400 text-sm py-10 text-center">Failed to load: {error}</div>}
        {!loading && !error && (
          <RidershipChart
            records={records}
            window={window}
            onWindowChange={setWindow}
            height={220}
          />
        )}
      </div>

      {/* Headway placeholder */}
      <div className="bg-gray-900 border border-dashed border-gray-700 rounded-lg px-5 py-4 opacity-60 flex items-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-gray-800 flex items-center justify-center shrink-0">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#6b7280" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
          </svg>
        </div>
        <div>
          <span className="text-[10px] bg-gray-800 text-gray-500 px-2 py-0.5 rounded-full">Coming soon</span>
          <p className="text-sm font-medium text-gray-400 mt-1">Headway Metrics</p>
          <p className="text-xs text-gray-600 mt-0.5">Track frequency, reliability, and schedule adherence for this route</p>
        </div>
      </div>
    </div>
  )
}
