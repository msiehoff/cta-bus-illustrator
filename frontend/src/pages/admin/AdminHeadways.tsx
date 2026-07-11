import { useMemo, useState } from 'react'
import { twMerge } from 'tailwind-merge'
import FilterValue from '../../components/admin/FilterValue'
import StatCard from '../../components/StatCard'
import { useHeadways } from '../../hooks/useHeadways'
import { useHeadwaySummary } from '../../hooks/useHeadwaySummary'
import { toChicagoServiceDate } from '../../lib/chicagoDate'

const formatTime = (value: string) => new Date(value).toLocaleString()

const formatMinutes = (mins: number) => {
  if (Number.isInteger(mins)) return String(mins)
  return mins.toFixed(1)
}

const formatCV = (cv: number) => (cv === 0 ? '—' : cv.toFixed(2))

const cvTrend = (cv: number) => {
  if (cv <= 0) return undefined
  if (cv < 0.3) return { text: 'Fairly regular', up: true as const }
  if (cv < 0.6) return { text: 'Moderate irregularity', up: undefined }
  return { text: 'Unreliable (bunches/gaps)', up: false as const }
}

const filterInputClass =
  'mt-1 block rounded-md bg-gray-950 border border-gray-800 px-3 py-2 text-white [color-scheme:dark]'

const AdminHeadways = () => {
  const [route, setRoute] = useState('')
  const [direction, setDirection] = useState('')
  const [stop, setStop] = useState('')
  const [vehicle, setVehicle] = useState('')
  const [date, setDate] = useState('')
  const [sortAsc, setSortAsc] = useState(false)
  const [offset, setOffset] = useState(0)
  const limit = 50

  const filters = { route, direction, stop, vehicle, date }

  const { data, loading, error } = useHeadways({
    ...filters,
    sort: sortAsc ? 'asc' : 'desc',
    limit,
    offset,
  })
  const { data: summary, error: summaryError } = useHeadwaySummary(filters)

  const page = useMemo(() => Math.floor(offset / limit) + 1, [offset, limit])
  const totalPages = data ? Math.max(1, Math.ceil(data.total / limit)) : 1

  const resetOffset = () => setOffset(0)

  const setFilter = (setter: (v: string) => void) => (value: string) => {
    setter(value)
    resetOffset()
  }

  const toggleTimeSort = () => {
    setSortAsc(prev => !prev)
    resetOffset()
  }

  const pooled = summary?.pooled
  const cvMeta = pooled ? cvTrend(pooled.cv) : undefined
  const showByStop = (summary?.byStop.length ?? 0) > 1
  const showEqualWeight = showByStop && summary != null

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold">Observed Headways</h2>
        <p className="text-sm text-gray-400 mt-1">
          Gaps between consecutive arrivals at a stop. Summary stats update with filters
          {summary?.source === 'stored' ? ' (from daily job)' : summary?.source === 'computed' ? ' (computed live)' : ''}.
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
          Route
          <input
            value={route}
            onChange={e => {
              setRoute(e.target.value)
              resetOffset()
            }}
            placeholder="e.g. 8"
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
            placeholder="Northbound"
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
        <label className="text-sm text-gray-400">
          Vehicle
          <input
            value={vehicle}
            onChange={e => {
              setVehicle(e.target.value)
              resetOffset()
            }}
            placeholder="from or to"
            className={twMerge(filterInputClass, 'w-32')}
          />
        </label>
        <div className="flex items-end">
          <button
            type="button"
            onClick={() => {
              setRoute('')
              setDirection('')
              setStop('')
              setVehicle('')
              resetOffset()
            }}
            disabled={!route && !direction && !stop && !vehicle}
            className={twMerge(
              'rounded-md border border-gray-800 px-3 py-2 text-sm',
              !route && !direction && !stop && !vehicle
                ? 'text-gray-600 cursor-not-allowed'
                : 'text-gray-300 hover:bg-gray-950 hover:text-white',
            )}
          >
            Clear filters
          </button>
        </div>
      </div>

      {(error || summaryError) && (
        <p className="text-red-400">{error || summaryError}</p>
      )}

      {pooled && pooled.count > 0 && (
        <div className="space-y-3">
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
            <StatCard label="Observations" value={String(pooled.count)} />
            <StatCard
              label="Mean headway"
              value={`${formatMinutes(pooled.meanMinutes)} min`}
              trend={
                showEqualWeight && summary
                  ? `Equal-stop avg ${formatMinutes(summary.equalStopWeight.meanMinutes)} min`
                  : undefined
              }
            />
            <StatCard
              label="Median headway"
              value={`${formatMinutes(pooled.medianMinutes)} min`}
            />
            <StatCard
              label="CV"
              value={formatCV(pooled.cv)}
              trend={cvMeta?.text}
              trendUp={cvMeta?.up}
            />
            <StatCard
              label="Avg wait ≈"
              value={`${formatMinutes(pooled.avgWaitMinutes)} min`}
              trend="mean / 2"
            />
          </div>

          {showByStop && summary && (
            <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
              <div className="px-4 py-3 border-b border-gray-800">
                <h3 className="text-sm font-medium text-gray-300">By stop</h3>
                <p className="text-xs text-gray-500 mt-0.5">
                  Per-stop means; equal-stop route average weights each stop the same.
                </p>
              </div>
              <div className="overflow-x-auto">
                <table className="min-w-full text-sm">
                  <thead className="bg-gray-950/80 text-gray-400">
                    <tr>
                      <th className="px-4 py-2 text-left font-medium">Stop</th>
                      <th className="px-4 py-2 text-left font-medium">n</th>
                      <th className="px-4 py-2 text-left font-medium">Mean</th>
                      <th className="px-4 py-2 text-left font-medium">Median</th>
                      <th className="px-4 py-2 text-left font-medium">CV</th>
                    </tr>
                  </thead>
                  <tbody>
                    {summary.byStop.map(s => (
                      <tr key={`${s.stopId}-${s.routeId}-${s.direction}`} className="border-t border-gray-800">
                        <td className="px-4 py-2">
                          <FilterValue
                            value={s.stopName || s.stopId}
                            active={stop === (s.stopName || s.stopId) || stop === s.stopId}
                            onSelect={setFilter(setStop)}
                            className="text-gray-300"
                          />
                        </td>
                        <td className="px-4 py-2 text-gray-400 tabular-nums">{s.count}</td>
                        <td className="px-4 py-2 text-white tabular-nums">{formatMinutes(s.meanMinutes)}</td>
                        <td className="px-4 py-2 text-gray-300 tabular-nums">{formatMinutes(s.medianMinutes)}</td>
                        <td className="px-4 py-2 text-gray-300 tabular-nums">{formatCV(s.cv)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}

      <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-950/80 text-gray-400">
              <tr>
                <th className="px-4 py-3 text-left font-medium">
                  <button
                    type="button"
                    onClick={toggleTimeSort}
                    className="inline-flex items-center gap-1 hover:text-white"
                  >
                    Time
                    <span className="text-xs text-gray-500">
                      {sortAsc ? '↑' : '↓'}
                    </span>
                  </button>
                </th>
                <th className="px-4 py-3 text-left font-medium">Minutes</th>
                <th className="px-4 py-3 text-left font-medium">Route</th>
                <th className="px-4 py-3 text-left font-medium">Direction</th>
                <th className="px-4 py-3 text-left font-medium">From → To</th>
                <th className="px-4 py-3 text-left font-medium">Stop</th>
              </tr>
            </thead>
            <tbody>
              {loading && !data ? (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                    Loading headways…
                  </td>
                </tr>
              ) : data?.headways.length ? (
                data.headways.map(h => {
                  const serviceDate = toChicagoServiceDate(h.timestamp)
                  const stopLabel = h.stopName || h.stopId
                  return (
                    <tr
                      key={`${h.stopId}-${h.routeId}-${h.direction}-${h.timestamp}-${h.toVehicleId}`}
                      className="border-t border-gray-800"
                    >
                      <td className="px-4 py-3 whitespace-nowrap">
                        <FilterValue
                          value={serviceDate}
                          active={date === serviceDate}
                          onSelect={setFilter(setDate)}
                          className="text-gray-300"
                        >
                          {formatTime(h.timestamp)}
                        </FilterValue>
                      </td>
                      <td className="px-4 py-3 text-white font-medium tabular-nums">
                        {formatMinutes(h.headwayMinutes)}
                      </td>
                      <td className="px-4 py-3">
                        <FilterValue
                          value={h.routeId}
                          active={route === h.routeId}
                          onSelect={setFilter(setRoute)}
                          className="text-white"
                        >
                          {h.routeName || h.routeId}
                        </FilterValue>
                      </td>
                      <td className="px-4 py-3">
                        <FilterValue
                          value={h.direction}
                          active={direction === h.direction}
                          onSelect={setFilter(setDirection)}
                          className="text-gray-300"
                        />
                      </td>
                      <td className="px-4 py-3 text-gray-300 whitespace-nowrap">
                        <FilterValue
                          value={h.fromVehicleId ?? ''}
                          active={vehicle === h.fromVehicleId}
                          onSelect={setFilter(setVehicle)}
                          className="text-gray-300"
                        />
                        <span className="mx-1 text-gray-600">→</span>
                        <FilterValue
                          value={h.toVehicleId ?? ''}
                          active={vehicle === h.toVehicleId}
                          onSelect={setFilter(setVehicle)}
                          className="text-gray-300"
                        />
                      </td>
                      <td className="px-4 py-3">
                        <FilterValue
                          value={stopLabel}
                          active={stop === stopLabel || stop === h.stopId}
                          onSelect={setFilter(setStop)}
                          className="text-gray-300"
                        />
                      </td>
                    </tr>
                  )
                })
              ) : (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                    No headways yet. Run a headway job first.
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

export default AdminHeadways
