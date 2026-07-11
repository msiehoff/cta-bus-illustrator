import { useMemo, useState } from 'react'
import { twMerge } from 'tailwind-merge'
import FilterValue from '../../components/admin/FilterValue'
import { useHeadwaySummaries } from '../../hooks/useHeadwaySummaries'

const formatMinutes = (mins: number) => {
  if (Number.isInteger(mins)) return String(mins)
  return mins.toFixed(1)
}

const formatCV = (cv: number) => (cv === 0 ? '—' : cv.toFixed(2))

const filterInputClass =
  'mt-1 block rounded-md bg-gray-950 border border-gray-800 px-3 py-2 text-white [color-scheme:dark]'

const GRAIN_OPTIONS = [
  { value: '', label: 'All grains' },
  { value: 'stop', label: 'stop' },
  { value: 'route_direction', label: 'route_direction' },
  { value: 'service_day', label: 'service_day' },
]

const METHOD_OPTIONS = [
  { value: '', label: 'All methods' },
  { value: 'pooled', label: 'pooled' },
  { value: 'equal_stop', label: 'equal_stop' },
]

const AdminHeadwaySummaries = () => {
  const [date, setDate] = useState('')
  const [grain, setGrain] = useState('stop')
  const [method, setMethod] = useState('')
  const [route, setRoute] = useState('')
  const [direction, setDirection] = useState('')
  const [stop, setStop] = useState('')
  const [sortAsc, setSortAsc] = useState(false)
  const [offset, setOffset] = useState(0)
  const limit = 50

  const { data, loading, error } = useHeadwaySummaries({
    date,
    grain,
    method,
    route,
    direction,
    stop,
    sort: sortAsc ? 'asc' : 'desc',
    limit,
    offset,
  })

  const page = useMemo(() => Math.floor(offset / limit) + 1, [offset, limit])
  const totalPages = data ? Math.max(1, Math.ceil(data.total / limit)) : 1

  const resetOffset = () => setOffset(0)

  const setFilter = (setter: (v: string) => void) => (value: string) => {
    setter(value)
    resetOffset()
  }

  const clearFilters = () => {
    setGrain('stop')
    setMethod('')
    setRoute('')
    setDirection('')
    setStop('')
    resetOffset()
  }

  const hasClearable = grain !== 'stop' || method || route || direction || stop

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold">Headway Summaries</h2>
        <p className="text-sm text-gray-400 mt-1">
          Persisted aggregates from the daily headway job. Defaults to stop-level rows.
        </p>
      </div>

      <div className="flex flex-wrap gap-3">
        <label className="text-sm text-gray-400">
          Service date
          <input
            type="date"
            value={date}
            onChange={e => {
              setDate(e.target.value)
              resetOffset()
            }}
            className={twMerge(filterInputClass, 'w-40')}
          />
        </label>
        <label className="text-sm text-gray-400">
          Grain
          <select
            value={grain}
            onChange={e => {
              setGrain(e.target.value)
              resetOffset()
            }}
            className={twMerge(filterInputClass, 'w-44')}
          >
            {GRAIN_OPTIONS.map(opt => (
              <option key={opt.value || 'all'} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </label>
        <label className="text-sm text-gray-400">
          Method
          <select
            value={method}
            onChange={e => {
              setMethod(e.target.value)
              resetOffset()
            }}
            className={twMerge(filterInputClass, 'w-40')}
          >
            {METHOD_OPTIONS.map(opt => (
              <option key={opt.value || 'all'} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </label>
        <label className="text-sm text-gray-400">
          Route
          <input
            value={route}
            onChange={e => {
              setRoute(e.target.value)
              resetOffset()
            }}
            placeholder="e.g. 66"
            className={twMerge(filterInputClass, 'w-28')}
          />
        </label>
        <label className="text-sm text-gray-400">
          Direction
          <input
            value={direction}
            onChange={e => {
              setDirection(e.target.value)
              resetOffset()
            }}
            placeholder="Westbound"
            className={twMerge(filterInputClass, 'w-40')}
          />
        </label>
        <label className="text-sm text-gray-400">
          Stop
          <input
            value={stop}
            onChange={e => {
              setStop(e.target.value)
              resetOffset()
            }}
            placeholder="Name or ID"
            className={twMerge(filterInputClass, 'w-48')}
          />
        </label>
        <div className="flex items-end">
          <button
            type="button"
            onClick={clearFilters}
            disabled={!hasClearable}
            className={twMerge(
              'rounded-md border border-gray-800 px-3 py-2 text-sm',
              !hasClearable
                ? 'text-gray-600 cursor-not-allowed'
                : 'text-gray-300 hover:bg-gray-950 hover:text-white',
            )}
          >
            Clear filters
          </button>
        </div>
      </div>

      {error && <p className="text-red-400">{error}</p>}

      <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-950/80 text-gray-400">
              <tr>
                <th className="px-4 py-3 text-left font-medium">
                  <button
                    type="button"
                    onClick={() => {
                      setSortAsc(prev => !prev)
                      resetOffset()
                    }}
                    className="inline-flex items-center gap-1 hover:text-white"
                  >
                    Date / Mean
                    <span className="text-xs text-gray-500">{sortAsc ? '↑' : '↓'}</span>
                  </button>
                </th>
                <th className="px-4 py-3 text-left font-medium">Grain</th>
                <th className="px-4 py-3 text-left font-medium">Method</th>
                <th className="px-4 py-3 text-left font-medium">Route</th>
                <th className="px-4 py-3 text-left font-medium">Direction</th>
                <th className="px-4 py-3 text-left font-medium">Stop</th>
                <th className="px-4 py-3 text-left font-medium">n</th>
                <th className="px-4 py-3 text-left font-medium">Mean</th>
                <th className="px-4 py-3 text-left font-medium">Median</th>
                <th className="px-4 py-3 text-left font-medium">CV</th>
                <th className="px-4 py-3 text-left font-medium">Wait≈</th>
              </tr>
            </thead>
            <tbody>
              {loading && !data ? (
                <tr>
                  <td colSpan={11} className="px-4 py-8 text-center text-gray-500">
                    Loading summaries…
                  </td>
                </tr>
              ) : data?.summaries.length ? (
                data.summaries.map(s => {
                  const stopLabel = s.stopName || s.stopId
                  return (
                    <tr
                      key={`${s.serviceDate}-${s.grain}-${s.method}-${s.routeId}-${s.direction}-${s.stopId}`}
                      className="border-t border-gray-800"
                    >
                      <td className="px-4 py-3 whitespace-nowrap">
                        <FilterValue
                          value={s.serviceDate}
                          active={date === s.serviceDate}
                          onSelect={setFilter(setDate)}
                          className="text-gray-300"
                        />
                      </td>
                      <td className="px-4 py-3">
                        <FilterValue
                          value={s.grain}
                          active={grain === s.grain}
                          onSelect={setFilter(setGrain)}
                          className="text-gray-300"
                        />
                      </td>
                      <td className="px-4 py-3">
                        <FilterValue
                          value={s.method}
                          active={method === s.method}
                          onSelect={setFilter(setMethod)}
                          className="text-gray-300"
                        />
                      </td>
                      <td className="px-4 py-3">
                        {s.routeId ? (
                          <FilterValue
                            value={s.routeId}
                            active={route === s.routeId}
                            onSelect={setFilter(setRoute)}
                            className="text-white"
                          >
                            {s.routeName || s.routeId}
                          </FilterValue>
                        ) : (
                          <span className="text-gray-600">—</span>
                        )}
                      </td>
                      <td className="px-4 py-3">
                        {s.direction ? (
                          <FilterValue
                            value={s.direction}
                            active={direction === s.direction}
                            onSelect={setFilter(setDirection)}
                            className="text-gray-300"
                          />
                        ) : (
                          <span className="text-gray-600">—</span>
                        )}
                      </td>
                      <td className="px-4 py-3">
                        {stopLabel ? (
                          <FilterValue
                            value={stopLabel}
                            active={stop === stopLabel || stop === s.stopId}
                            onSelect={setFilter(setStop)}
                            className="text-gray-300"
                          />
                        ) : (
                          <span className="text-gray-600">—</span>
                        )}
                      </td>
                      <td className="px-4 py-3 text-gray-400 tabular-nums">{s.count}</td>
                      <td className="px-4 py-3 text-white font-medium tabular-nums">
                        {formatMinutes(s.meanMinutes)}
                      </td>
                      <td className="px-4 py-3 text-gray-300 tabular-nums">
                        {formatMinutes(s.medianMinutes)}
                      </td>
                      <td className="px-4 py-3 text-gray-300 tabular-nums">{formatCV(s.cv)}</td>
                      <td className="px-4 py-3 text-gray-400 tabular-nums">
                        {formatMinutes(s.avgWaitMinutes)}
                      </td>
                    </tr>
                  )
                })
              ) : (
                <tr>
                  <td colSpan={11} className="px-4 py-8 text-center text-gray-500">
                    No summaries yet. Run a headway job first.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="flex items-center justify-between px-4 py-3 border-t border-gray-800 text-sm text-gray-400">
          <span>{data ? `${data.total} total` : '—'}</span>
          <div className="flex items-center gap-2">
            <button
              type="button"
              disabled={offset === 0}
              onClick={() => setOffset(Math.max(0, offset - limit))}
              className={twMerge(
                'px-3 py-1 rounded-md border border-gray-800',
                offset === 0 ? 'opacity-40 cursor-not-allowed' : 'hover:bg-gray-950',
              )}
            >
              Previous
            </button>
            <span>
              Page {page} / {totalPages}
            </span>
            <button
              type="button"
              disabled={!data || offset + limit >= data.total}
              onClick={() => setOffset(offset + limit)}
              className={twMerge(
                'px-3 py-1 rounded-md border border-gray-800',
                !data || offset + limit >= data.total
                  ? 'opacity-40 cursor-not-allowed'
                  : 'hover:bg-gray-950',
              )}
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default AdminHeadwaySummaries
