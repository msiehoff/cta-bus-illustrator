import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useSystemRidership } from '../hooks/useSystemRidership'
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

export default function SystemOverview() {
  const { records, loading, error } = useSystemRidership()
  const navigate = useNavigate()
  const [search, setSearch] = useState('')
  const [window, setWindow] = useState<WindowKey>('5y')
  const [routes, setRoutes] = useState<{ routeId: string; routeName: string; avgRides: number }[]>([])

  useEffect(() => {
    fetch('/api/v1/routes')
      .then(r => r.json())
      .then(data => {
        const sorted = (data.features ?? [])
          .filter((f: { properties: { avgRides?: number } }) => f.properties.avgRides != null)
          .map((f: { properties: { routeId: string; routeName: string; avgRides: number } }) => ({
            routeId: f.properties.routeId,
            routeName: f.properties.routeName,
            avgRides: f.properties.avgRides,
          }))
          .sort((a: { avgRides: number }, b: { avgRides: number }) => b.avgRides - a.avgRides)
        setRoutes(sorted)
      })
      .catch(() => {})
  }, [])

  // Latest month in the dataset — anchor for window cutoffs
  const latestMonth = useMemo(() => {
    if (!records.length) return null
    return records[records.length - 1].month
  }, [records])

  // Latest month values for stat cards
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

  // First month of the current window — baseline for % comparison
  const windowStart = useMemo(() => {
    if (!records.length || !latestMonth) return null
    const cutoff = cutoffMonth(window, latestMonth)
    const windowRecords = cutoff ? records.filter(r => r.month >= cutoff) : records
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

  const filteredRoutes = routes.filter(r =>
    r.routeName.toLowerCase().includes(search.toLowerCase()) ||
    r.routeId.includes(search)
  )

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">System Overview</h1>
        <p className="text-sm text-gray-400 mt-0.5">
          CTA Bus Network · All Routes{latest ? ` · ${formatMonth(latest.month)}` : ''}
        </p>
      </div>

      {/* Stat cards */}
      {latest && (
        <div className="grid grid-cols-3 gap-3 mb-5">
          <StatCard label={`Total weekday riders · ${formatMonth(latest.month)}`} {...statProps('weekday')} />
          <StatCard label={`Total Saturday riders · ${formatMonth(latest.month)}`} {...statProps('saturday')} />
          <StatCard label={`Total Sunday riders · ${formatMonth(latest.month)}`} {...statProps('sunday')} />
        </div>
      )}

      {/* Chart — window owned here so buttons work and stat cards stay in sync */}
      <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4 mb-5">
        <h2 className="text-sm font-medium text-white mb-4">System ridership over time</h2>
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

      {/* Routes table */}
      <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
        <div className="flex items-center justify-between px-5 py-3 border-b border-gray-800">
          <h2 className="text-sm font-medium text-white">Routes</h2>
          <input
            type="text"
            placeholder="Search routes…"
            value={search}
            onChange={e => setSearch(e.target.value)}
            className="bg-gray-800 border border-gray-700 rounded-md px-3 py-1.5 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-gray-500 w-44"
          />
        </div>
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-gray-800/50 text-xs text-gray-500">
              <th className="text-left px-5 py-2 font-normal w-10">#</th>
              <th className="text-left px-3 py-2 font-normal">Route</th>
              <th className="text-right px-5 py-2 font-normal">Avg weekday riders</th>
            </tr>
          </thead>
          <tbody>
            {filteredRoutes.map((route, i) => (
              <tr
                key={route.routeId}
                onClick={() => navigate(`/routes/${route.routeId}`, { state: { routeName: route.routeName } })}
                className="border-t border-gray-800 hover:bg-gray-800/40 cursor-pointer transition-colors"
              >
                <td className="px-5 py-2.5 text-gray-600 text-xs">{i + 1}</td>
                <td className="px-3 py-2.5">
                  <div className="flex items-center gap-2.5">
                    <span className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-blue-900 text-blue-300 text-[10px] font-semibold shrink-0">
                      {route.routeId}
                    </span>
                    <span className="text-white">{route.routeName}</span>
                  </div>
                </td>
                <td className="px-5 py-2.5 text-right text-gray-300">
                  {formatRides(route.avgRides)}
                </td>
              </tr>
            ))}
            {filteredRoutes.length === 0 && (
              <tr>
                <td colSpan={3} className="px-5 py-8 text-center text-gray-600 text-xs">
                  No routes found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
