import { useMemo, useState } from 'react'
import { twMerge } from 'tailwind-merge'
import { useArrivals } from '../../hooks/useArrivals'

const formatTime = (value: string) => new Date(value).toLocaleString()

const AdminArrivals = () => {
  const [route, setRoute] = useState('')
  const [direction, setDirection] = useState('')
  const [offset, setOffset] = useState(0)
  const limit = 50

  const { data, loading, error } = useArrivals({ route, direction, limit, offset })

  const page = useMemo(() => Math.floor(offset / limit) + 1, [offset, limit])
  const totalPages = data ? Math.max(1, Math.ceil(data.total / limit)) : 1

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold">Recent Arrivals</h2>
        <p className="text-sm text-gray-400 mt-1">
          Detected stop arrivals from the pipeline. Refreshes every 10 seconds.
        </p>
      </div>

      <div className="flex flex-wrap gap-3">
        <label className="text-sm text-gray-400">
          Route
          <input
            value={route}
            onChange={e => {
              setRoute(e.target.value)
              setOffset(0)
            }}
            placeholder="e.g. 8"
            className="mt-1 block w-28 rounded-md bg-gray-950 border border-gray-800 px-3 py-2 text-white"
          />
        </label>
        <label className="text-sm text-gray-400">
          Direction
          <input
            value={direction}
            onChange={e => {
              setDirection(e.target.value)
              setOffset(0)
            }}
            placeholder="Northbound"
            className="mt-1 block w-40 rounded-md bg-gray-950 border border-gray-800 px-3 py-2 text-white"
          />
        </label>
      </div>

      {error && <p className="text-red-400">{error}</p>}

      <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-950/80 text-gray-400">
              <tr>
                <th className="px-4 py-3 text-left font-medium">Time</th>
                <th className="px-4 py-3 text-left font-medium">Route</th>
                <th className="px-4 py-3 text-left font-medium">Direction</th>
                <th className="px-4 py-3 text-left font-medium">Vehicle</th>
                <th className="px-4 py-3 text-left font-medium">Stop</th>
              </tr>
            </thead>
            <tbody>
              {loading && !data ? (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-gray-500">
                    Loading arrivals…
                  </td>
                </tr>
              ) : data?.arrivals.length ? (
                data.arrivals.map(arrival => (
                  <tr key={`${arrival.vehicleId}-${arrival.timestamp}-${arrival.stopId}`} className="border-t border-gray-800">
                    <td className="px-4 py-3 text-gray-300 whitespace-nowrap">{formatTime(arrival.timestamp)}</td>
                    <td className="px-4 py-3 text-white">{arrival.routeId}</td>
                    <td className="px-4 py-3 text-gray-300">{arrival.direction}</td>
                    <td className="px-4 py-3 text-gray-300">{arrival.vehicleId}</td>
                    <td className="px-4 py-3 text-gray-300">{arrival.stopId}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-gray-500">
                    No arrivals yet.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="flex items-center justify-between px-4 py-3 border-t border-gray-800 text-sm text-gray-400">
          <span>
            {data ? `${data.total} total` : '—'}
          </span>
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
            <span>Page {page} / {totalPages}</span>
            <button
              type="button"
              disabled={!data || offset + limit >= data.total}
              onClick={() => setOffset(offset + limit)}
              className={twMerge(
                'px-3 py-1 rounded-md border border-gray-800',
                !data || offset + limit >= data.total ? 'opacity-40 cursor-not-allowed' : 'hover:bg-gray-950',
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

export default AdminArrivals
