import { useMemo, useState, type ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'
import { twMerge } from 'tailwind-merge'
import type { HeadwayRoutePeriod } from '../types/api'
import {
  formatHeadwayCV,
  formatHeadwayMinutes,
  HEADWAY_CONSISTENCY_LABEL,
  HEADWAY_CONSISTENCY_TOOLTIP,
} from '../lib/headwayUtils'
import HintLabel from './HintLabel'

type SortColumn = 'route' | 'median' | 'wait' | 'cv'

interface SortState {
  column: SortColumn
  dir: 'asc' | 'desc'
}

interface Props {
  routes: HeadwayRoutePeriod[]
  loading?: boolean
  search?: string
  onSearchChange?: (value: string) => void
}

const HeadwayRoutesTable = ({
  routes,
  loading = false,
  search = '',
  onSearchChange,
}: Props) => {
  const navigate = useNavigate()
  const [sort, setSort] = useState<SortState>({ column: 'median', dir: 'desc' })

  const toggleSort = (column: SortColumn) =>
    setSort(prev =>
      prev.column === column
        ? { column, dir: prev.dir === 'desc' ? 'asc' : 'desc' }
        : { column, dir: column === 'route' ? 'asc' : 'desc' },
    )

  const filtered = useMemo(() => {
    const q = search.toLowerCase()
    return routes.filter(
      r =>
        (r.routeName ?? '').toLowerCase().includes(q) ||
        r.routeId.toLowerCase().includes(q),
    )
  }, [routes, search])

  const sorted = useMemo(() => {
    const copy = [...filtered]
    const { column, dir } = sort
    const factor = dir === 'asc' ? 1 : -1
    copy.sort((a, b) => {
      switch (column) {
        case 'route':
          return (a.routeName || a.routeId).localeCompare(b.routeName || b.routeId) * factor
        case 'median':
          return (a.medianMinutes - b.medianMinutes) * factor
        case 'wait':
          return (a.avgWaitMinutes - b.avgWaitMinutes) * factor
        case 'cv':
          return (a.cv - b.cv) * factor
      }
    })
    return copy
  }, [filtered, sort])

  const SortableHeader = ({
    column,
    align = 'left',
    children,
    hint,
  }: {
    column: SortColumn
    align?: 'left' | 'right'
    children: ReactNode
    hint?: string
  }) => {
    const active = sort.column === column
    const sortMark = (
      <span className="text-[10px] text-gray-600">
        {active ? (sort.dir === 'asc' ? '↑' : '↓') : ''}
      </span>
    )

    if (hint) {
      return (
        <th
          className={twMerge(
            'px-3 py-2 font-normal relative',
            align === 'right' && 'text-right',
          )}
        >
          <div
            className={twMerge(
              'inline-flex items-center gap-1',
              align === 'right' && 'justify-end w-full',
              active ? 'text-white' : 'text-gray-500',
            )}
          >
            <HintLabel hint={hint} align={align === 'right' ? 'right' : 'left'} placement="above">
              {children}
            </HintLabel>
            <button
              type="button"
              onClick={() => toggleSort(column)}
              className="hover:text-white transition-colors"
              aria-label={`Sort by ${column}`}
            >
              {sortMark}
            </button>
          </div>
        </th>
      )
    }

    return (
      <th
        className={twMerge(
          'px-3 py-2 font-normal',
          align === 'right' && 'text-right',
        )}
      >
        <button
          type="button"
          onClick={() => toggleSort(column)}
          className={twMerge(
            'inline-flex items-center gap-1 hover:text-white transition-colors',
            active ? 'text-white' : 'text-gray-500',
            align === 'right' && 'justify-end w-full',
          )}
        >
          {children}
          {sortMark}
        </button>
      </th>
    )
  }

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
      {onSearchChange && (
        <div className="px-3 py-2 border-b border-gray-800">
          <input
            type="search"
            value={search}
            onChange={e => onSearchChange(e.target.value)}
            placeholder="Search routes…"
            className="w-full bg-gray-950 border border-gray-800 rounded px-3 py-1.5 text-sm text-white placeholder:text-gray-600 focus:outline-none focus:border-gray-600"
          />
        </div>
      )}

      <div className="overflow-x-auto">
        <table className="w-full text-sm text-left">
          <thead className="text-xs text-gray-500 border-b border-gray-800">
            <tr>
              <SortableHeader column="route">Route</SortableHeader>
              <SortableHeader column="median" align="right">Median headway</SortableHeader>
              <SortableHeader column="wait" align="right">Avg wait</SortableHeader>
              <SortableHeader
                column="cv"
                align="right"
                hint={HEADWAY_CONSISTENCY_TOOLTIP}
              >
                {HEADWAY_CONSISTENCY_LABEL}
              </SortableHeader>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr>
                <td colSpan={4} className="px-3 py-10 text-center text-gray-500">
                  Loading…
                </td>
              </tr>
            )}
            {!loading && sorted.length === 0 && (
              <tr>
                <td colSpan={4} className="px-3 py-10 text-center text-gray-500">
                  No headway summaries yet. Run a daily headway job to populate this table.
                </td>
              </tr>
            )}
            {!loading &&
              sorted.map(route => (
                <tr
                  key={route.routeId}
                  onClick={() => navigate(`/headways/routes/${encodeURIComponent(route.routeId)}`)}
                  className="border-b border-gray-800/80 hover:bg-gray-800/40 cursor-pointer transition-colors"
                >
                  <td className="px-3 py-2.5">
                    <span className="text-white font-medium">{route.routeId}</span>
                    {route.routeName && (
                      <span className="text-gray-500 ml-2">{route.routeName}</span>
                    )}
                  </td>
                  <td className="px-3 py-2.5 text-right text-white tabular-nums">
                    {formatHeadwayMinutes(route.medianMinutes)} min
                  </td>
                  <td className="px-3 py-2.5 text-right text-gray-300 tabular-nums">
                    {formatHeadwayMinutes(route.avgWaitMinutes)} min
                  </td>
                  <td className="px-3 py-2.5 text-right text-gray-400 tabular-nums">
                    {formatHeadwayCV(route.cv)}
                  </td>
                </tr>
              ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default HeadwayRoutesTable
