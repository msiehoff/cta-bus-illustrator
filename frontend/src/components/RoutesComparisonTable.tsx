import { useMemo, useState, type ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'
import { twMerge } from 'tailwind-merge'
import {
  formatPct,
  formatRecoveryPct,
  formatRides,
  recoveryColorClass,
} from '../lib/ridershipUtils'
import type { RouteComparison } from '../types/api'

type SortColumn = 'route' | 'current' | 'recovery' | 'yearAgo' | 'fiveYear'

interface SortState {
  column: SortColumn
  dir: 'asc' | 'desc'
}

interface Props {
  routes: RouteComparison[]
  loading?: boolean
  showHeader?: boolean
  search?: string
  onSearchChange?: (value: string) => void
}

const compareNullable = (
  a: number | undefined,
  b: number | undefined,
  dir: 'asc' | 'desc',
): number => {
  if (a == null && b == null) return 0
  if (a == null) return 1
  if (b == null) return -1
  return dir === 'asc' ? a - b : b - a
}

const RoutesComparisonTable = ({
  routes,
  loading = false,
  showHeader = true,
  search = '',
  onSearchChange,
}: Props) => {
  const navigate = useNavigate()
  const [sort, setSort] = useState<SortState>({ column: 'current', dir: 'desc' })

  const toggleSort = (column: SortColumn) =>
    setSort(prev =>
      prev.column === column
        ? { column, dir: prev.dir === 'desc' ? 'asc' : 'desc' }
        : { column, dir: column === 'route' ? 'asc' : 'desc' },
    )

  const filtered = useMemo(
    () => routes.filter(r =>
      r.routeName.toLowerCase().includes(search.toLowerCase()) ||
      r.routeId.includes(search),
    ),
    [routes, search],
  )

  const sorted = useMemo(() => {
    const copy = [...filtered]
    const { column, dir } = sort
    const factor = dir === 'asc' ? 1 : -1

    copy.sort((a, b) => {
      switch (column) {
        case 'route':
          return a.routeName.localeCompare(b.routeName) * factor
        case 'current':
          return (a.current - b.current) * factor
        case 'recovery':
          return compareNullable(a.recoveryPct, b.recoveryPct, dir)
        case 'yearAgo':
          return compareNullable(a.yearAgoPct, b.yearAgoPct, dir)
        case 'fiveYear':
          return compareNullable(a.fiveYearPct, b.fiveYearPct, dir)
      }
    })
    return copy
  }, [filtered, sort])

  const formatDelta = (pct: number | undefined) => {
    if (pct == null) return <span className="text-gray-600">—</span>
    const up = pct >= 0
    return (
      <span className={up ? 'text-green-400' : 'text-red-400'}>
        {formatPct(pct)}
      </span>
    )
  }

  const SortableHeader = ({
    column,
    align = 'left',
    className,
    children,
  }: {
    column: SortColumn
    align?: 'left' | 'right'
    className?: string
    children: ReactNode
  }) => {
    const active = sort.column === column
    return (
      <th className={twMerge(
        'px-3 py-2 font-normal',
        align === 'right' && 'text-right',
        className,
      )}>
        <button
          type="button"
          onClick={() => toggleSort(column)}
          className={twMerge(
            'inline-flex items-center gap-1 transition-colors hover:text-gray-300',
            align === 'right' && 'flex-row-reverse',
            active ? 'text-gray-300' : 'text-gray-500',
          )}
        >
          {children}
          {active && (
            <span className="text-[10px] text-gray-400" aria-hidden>
              {sort.dir === 'desc' ? '↓' : '↑'}
            </span>
          )}
        </button>
      </th>
    )
  }

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
      {showHeader && (
        <div className="px-5 py-3 border-b border-gray-800 flex items-center justify-between gap-3 flex-wrap">
          <span className="text-xs text-gray-500">{filtered.length} routes</span>
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
              <tr className="bg-gray-800/50 text-xs">
                <th className="text-left px-5 py-2 font-normal text-gray-500 w-10">#</th>
                <SortableHeader column="route" className="text-left px-3">
                  Route
                </SortableHeader>
                <SortableHeader column="current" align="right">
                  Now
                </SortableHeader>
                <SortableHeader column="recovery" align="right">
                  vs 2019
                </SortableHeader>
                <SortableHeader column="yearAgo" align="right" className="hidden sm:table-cell">
                  1y
                </SortableHeader>
                <SortableHeader column="fiveYear" align="right" className="hidden sm:table-cell pr-5">
                  5y
                </SortableHeader>
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
