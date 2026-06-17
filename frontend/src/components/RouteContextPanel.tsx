import { Link } from 'react-router-dom'
import type { RouteComparison } from '../types/api'
import {
  formatRecoveryPct,
  formatRides,
  getPeerRoutes,
  getRouteRank,
  recoveryColorClass,
} from '../lib/ridershipUtils'
import { getPairedRouteId } from '../lib/expressPairs'

interface Props {
  routeId: string
  routeName: string
  comparison: RouteComparison[]
  systemTotal: number
  ridershipType: string
  selectedMonth: string
}

const RouteContextPanel = ({
  routeId,
  routeName,
  comparison,
  systemTotal,
  ridershipType,
  selectedMonth,
}: Props) => {
  const rank = getRouteRank(routeId, comparison, systemTotal)
  const peers = getPeerRoutes(routeId, comparison)
  const pairedId = getPairedRouteId(routeId)
  const pairedRoute = pairedId ? comparison.find(r => r.routeId === pairedId) : undefined
  const currentRoute = comparison.find(r => r.routeId === routeId)

  const mapLink = `/?route=${encodeURIComponent(routeId)}&month=${selectedMonth}&type=${ridershipType}`

  return (
    <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-3 mb-5">
      {rank && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
          <p className="text-xs text-gray-500 mb-1">Network rank</p>
          <p className="text-2xl font-semibold text-white">#{rank.rank}</p>
          <p className="text-xs text-gray-500 mt-1">of {rank.total} routes</p>
        </div>
      )}
      {rank && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
          <p className="text-xs text-gray-500 mb-1">Network share</p>
          <p className="text-2xl font-semibold text-white">{rank.networkSharePct.toFixed(1)}%</p>
          <p className="text-xs text-gray-500 mt-1">of total {ridershipType} riders</p>
        </div>
      )}
      <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3 sm:col-span-2 lg:col-span-1">
        <p className="text-xs text-gray-500 mb-1">View on map</p>
        <Link
          to={mapLink}
          className="inline-flex items-center gap-1.5 text-sm text-red-400 hover:text-red-300 transition-colors mt-1"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <polygon points="3 6 9 3 15 6 21 3 21 18 15 21 9 18 3 21" />
          </svg>
          Show {routeId} on map
        </Link>
      </div>

      {pairedRoute && currentRoute && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3 sm:col-span-2 lg:col-span-4">
          <p className="text-xs text-gray-500 uppercase tracking-wide mb-2">Corridor (local + express)</p>
          <div className="flex flex-wrap items-center gap-4 text-sm">
            <div>
              <span className="text-gray-400">{routeName}</span>
              <span className="text-white font-medium ml-2">{formatRides(currentRoute.current)}</span>
            </div>
            <span className="text-gray-600">+</span>
            <div>
              <span className="text-gray-400">{pairedRoute.routeName}</span>
              <span className="text-white font-medium ml-2">{formatRides(pairedRoute.current)}</span>
            </div>
            <span className="text-gray-600">=</span>
            <div>
              <span className="text-gray-400">Combined</span>
              <span className="text-white font-semibold ml-2">
                {formatRides(currentRoute.current + pairedRoute.current)}
              </span>
            </div>
            {currentRoute.recoveryPct != null && pairedRoute.recoveryPct != null && (
              <span className={`text-xs ${recoveryColorClass((currentRoute.recoveryPct + pairedRoute.recoveryPct) / 2)}`}>
                {formatRecoveryPct((currentRoute.recoveryPct + pairedRoute.recoveryPct) / 2)} avg vs pre-pandemic
              </span>
            )}
          </div>
        </div>
      )}

      {peers.length > 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3 sm:col-span-2 lg:col-span-4">
          <p className="text-xs text-gray-500 uppercase tracking-wide mb-2">Similar ridership</p>
          <div className="flex flex-wrap gap-2">
            {peers.map(peer => (
              <Link
                key={peer.routeId}
                to={`/routes/${peer.routeId}`}
                state={{ routeName: peer.routeName }}
                className="inline-flex items-center gap-2 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded-lg px-3 py-1.5 text-xs transition-colors"
              >
                <span className="text-white truncate max-w-[10rem]">{peer.routeName}</span>
                <span className="text-gray-300">{formatRides(peer.current)}</span>
                {peer.recoveryPct != null && (
                  <span className={recoveryColorClass(peer.recoveryPct)}>
                    {formatRecoveryPct(peer.recoveryPct)}
                  </span>
                )}
              </Link>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default RouteContextPanel
