import { useMemo } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useRouteRidership } from '../hooks/useRouteRidership'
import RidershipChart from '../components/RidershipChart'
import StatCard from '../components/StatCard'

function formatRides(value: number): string {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(1)}k`
  return value.toFixed(0)
}

export default function RoutePage() {
  const { externalId = '' } = useParams<{ externalId: string }>()
  const navigate = useNavigate()
  const { records, loading, error } = useRouteRidership(externalId)

  // Compute latest-month averages for the stat cards
  const stats = useMemo(() => {
    if (!records.length) return null
    const latestMonth = records[records.length - 1].month
    const latest = records.filter(r => r.month === latestMonth)
    const get = (type: string) => latest.find(r => r.type === type)?.avgRides ?? 0
    return {
      weekday: get('weekday'),
      saturday: get('saturday'),
      sunday: get('sunday'),
    }
  }, [records])

  // Try to get a route name from the first record (not available from ridership endpoint,
  // so fall back to showing the ID). You could also pass routeName via router state.
  const routeName = (history.state as { routeName?: string } | null)?.routeName ?? externalId

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
          {stats && (
            <p className="text-sm text-gray-400 mt-0.5">
              {formatRides(stats.weekday)} avg weekday riders
            </p>
          )}
        </div>
      </div>

      {/* Stat cards */}
      {stats && (
        <div className="grid grid-cols-3 gap-3 mb-5">
          <StatCard label="Avg Weekday Riders" value={formatRides(stats.weekday)} />
          <StatCard label="Avg Saturday Riders" value={formatRides(stats.saturday)} />
          <StatCard label="Avg Sunday Riders" value={formatRides(stats.sunday)} />
        </div>
      )}

      {/* Chart */}
      <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-4">Ridership over time</h2>
        {loading && <div className="text-gray-500 text-sm py-10 text-center">Loading…</div>}
        {error && <div className="text-red-400 text-sm py-10 text-center">Failed to load: {error}</div>}
        {!loading && !error && <RidershipChart records={records} height={220} />}
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
