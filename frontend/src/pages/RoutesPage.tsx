import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

interface RouteEntry {
  routeId: string
  routeName: string
  avgRides: number
}

function formatRides(value: number): string {
  if (value >= 1_000) return `${(value / 1_000).toFixed(1)}k`
  return value.toFixed(0)
}

export default function RoutesPage() {
  const navigate = useNavigate()
  const [routes, setRoutes] = useState<RouteEntry[]>([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/routes')
      .then(r => r.json())
      .then(data => {
        const sorted: RouteEntry[] = (data.features ?? [])
          .filter((f: { properties: { avgRides?: number } }) => f.properties.avgRides != null)
          .map((f: { properties: { routeId: string; routeName: string; avgRides: number } }) => ({
            routeId: f.properties.routeId,
            routeName: f.properties.routeName,
            avgRides: f.properties.avgRides,
          }))
          .sort((a: RouteEntry, b: RouteEntry) => b.avgRides - a.avgRides)
        setRoutes(sorted)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const filtered = routes.filter(r =>
    r.routeName.toLowerCase().includes(search.toLowerCase()) ||
    r.routeId.includes(search)
  )

  return (
    <div className="p-6 max-w-3xl mx-auto">
      <div className="mb-5">
        <h1 className="text-xl font-semibold text-white">Routes</h1>
        <p className="text-sm text-gray-400 mt-0.5">All CTA bus routes · sorted by ridership</p>
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
        <div className="px-5 py-3 border-b border-gray-800 flex items-center justify-between">
          <span className="text-xs text-gray-500">{filtered.length} routes</span>
          <input
            type="text"
            placeholder="Search routes…"
            value={search}
            onChange={e => setSearch(e.target.value)}
            className="bg-gray-800 border border-gray-700 rounded-md px-3 py-1.5 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-gray-500 w-44"
          />
        </div>
        {loading ? (
          <div className="py-12 text-center text-gray-500 text-sm">Loading…</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-gray-800/50 text-xs text-gray-500">
                <th className="text-left px-5 py-2 font-normal w-10">#</th>
                <th className="text-left px-3 py-2 font-normal">Route</th>
                <th className="text-right px-5 py-2 font-normal">Avg Weekday Riders</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((route, i) => (
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
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
