import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  formatPct,
  formatRecoveryPct,
  formatRides,
  recoveryColorClass,
} from '../lib/ridershipUtils'
import type { RouteComparison } from '../types/api'

type SortKey = 'ridership' | 'recovery'

interface Props {
  routes: RouteComparison[]
  loading?: boolean
  showHeader?: boolean
  search?: string
  onSearchChange?: (value: string) => void
}

const RoutesComparisonTable = ({
  routes,
  loading = false,
  showHeader = true,
  search = '',
  onSearchChange,
}: Props) => {
  const navigate = useNavigate()
  const [sortBy, setSortBy] = useState<SortKey>('ridership')

  const filtered = useMemo(
    () => routes.filter(r =>
      r.routeName.toLowerCase().includes(search.toLowerCase()) ||
      r.routeId.includes(search),
    ),
    [routes, search],
  )

  const sorted = useMemo(() => {
    const copy = [...filtered]
    if (sortBy === 'recovery') {
      copy.sort((a, b) => (b.recoveryPct ?? -1) - (a.recoveryPct ?? -1))
    } else {
      copy.sort((a, b) => b.current - a.current)
    }
    return copy
  }, [filtered, sortBy])

  const formatDelta = (pct: number | undefined) => {
    if (pct == null) return <span className="text-gray-600">—</span>
    const up = pct >= 0
    return (
      <span className={up ? 'text-green-400' : 'text-red-400'}>
        {formatPct(pct)}
      </span>
    )
  }

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
      {showHeader && (
        <div className="px-5 py-3 border-b border-gray-800 flex items-center justify-between gap-3 flex-wrap">
          <div className="flex items-center gap-3">
            <span className="text-xs text-gray-500">{filtered.length} routes</span>
            <div className="flex items-center gap-1">
              <span className="text-[10px] text-gray-600">Sort:</span>
              {(['ridership', 'recovery'] as const).map(key => (
                <button
                  key={key}
                  onClick={() => setSortBy(key)}
                  className={`px-2 py-0.5 rounded text-[10px] transition-colors ${
                    sortBy === key
                      ? 'bg-gray-700 text-white'
                      : 'text-gray-500 hover:text-gray-300'
                  }`}
                >
                  {key === 'ridership' ? 'Ridership' : 'Recovery'}
                </button>
              ))}
            </div>
          </div>
          {onSearchChange && (
            <input
              type="text"
              placeholder="Search routes…"
              value={search}
              onChange={e => onSearchChange(e.target.value)}
              className="bg-gray-800 border border-gray-700 rounded-md px-3 py-1.5 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-gray-500 w-44"
            />
          )}
        </div>
      )}

      {loading ? (
        <div className="py-12 text-center text-gray-500 text-sm">Loading…</div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm min-w-[32rem]">
            <thead>
              <tr className="bg-gray-800/50 text-xs text-gray-500">
                <th className="text-left px-5 py-2 font-normal w-10">#</th>
                <th className="text-left px-3 py-2 font-normal">Route</th>
                <th className="text-right px-3 py-2 font-normal">Now</th>
                <th className="text-right px-3 py-2 font-normal">vs 2019</th>
                <th className="text-right px-3 py-2 font-normal hidden sm:table-cell">1y</th>
                <th className="text-right px-5 py-2 font-normal hidden sm:table-cell">5y</th>
              </tr>
            </thead>
            <tbody>
              {sorted.map((route, i) => (
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
                  <td className="px-3 py-2.5 text-right text-gray-300">
                    {formatRides(route.current)}
                  </td>
                  <td className="px-3 py-2.5 text-right">
                    {route.recoveryPct != null ? (
                      <span className={recoveryColorClass(route.recoveryPct)}>
                        {formatRecoveryPct(route.recoveryPct)}
                      </span>
                    ) : (
                      <span className="text-gray-600">—</span>
                    )}
                  </td>
                  <td className="px-3 py-2.5 text-right hidden sm:table-cell">
                    {formatDelta(route.yearAgoPct)}
                  </td>
                  <td className="px-5 py-2.5 text-right hidden sm:table-cell">
                    {formatDelta(route.fiveYearPct)}
                  </td>
                </tr>
              ))}
              {sorted.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-5 py-8 text-center text-gray-600 text-xs">
                    No routes found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

export default RoutesComparisonTable
