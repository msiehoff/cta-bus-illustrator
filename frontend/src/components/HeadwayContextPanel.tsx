import {
  formatHeadwayMinutes,
  getHeadwayRank,
  headwayVsNetworkPct,
} from '../lib/headwayUtils'
import type { HeadwayRoutePeriod } from '../types/api'

interface Props {
  routeId: string
  routes: HeadwayRoutePeriod[]
  routeMedian: number
  routeWait: number
  networkMedian: number
  networkWait: number
  periodLabel?: string | null
}

const formatDelta = (pct: number | null, unitLabel: string) => {
  if (pct == null) return null
  const better = pct < 0
  const abs = Math.abs(pct)
  const label = abs < 0.5 ? 'about even with' : `${abs.toFixed(0)}% ${better ? 'shorter than' : 'longer than'}`
  return {
    text: `${label} network ${unitLabel}`,
    better: abs < 0.5 ? undefined : better,
  }
}

const HeadwayContextPanel = ({
  routeId,
  routes,
  routeMedian,
  routeWait,
  networkMedian,
  networkWait,
  periodLabel,
}: Props) => {
  const rank = getHeadwayRank(routeId, routes)
  const medianDelta = formatDelta(headwayVsNetworkPct(routeMedian, networkMedian), 'median')
  const waitDelta = formatDelta(headwayVsNetworkPct(routeWait, networkWait), 'wait')

  if (!rank && !medianDelta) return null

  return (
    <div className="mb-5">
      {periodLabel && (
        <p className="text-xs text-gray-500 mb-2">
          Network context · {periodLabel}
        </p>
      )}
      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {rank && (
          <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
            <p className="text-xs text-gray-500 mb-1">Network rank</p>
            <p className="text-2xl font-semibold text-white">#{rank.rank}</p>
            <p className="text-xs text-gray-500 mt-1">
              of {rank.total} routes by median headway
              <span className="text-gray-600"> · #1 = shortest</span>
            </p>
          </div>
        )}
        {medianDelta && (
          <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
            <p className="text-xs text-gray-500 mb-1">Vs network median</p>
            <p className="text-2xl font-semibold text-white">
              {formatHeadwayMinutes(routeMedian)} min
            </p>
            <p
              className={`text-xs mt-1 ${
                medianDelta.better === true
                  ? 'text-green-400'
                  : medianDelta.better === false
                    ? 'text-red-400'
                    : 'text-gray-500'
              }`}
            >
              {medianDelta.text}
              {' · '}
              network {formatHeadwayMinutes(networkMedian)} min
            </p>
          </div>
        )}
        {waitDelta && (
          <div className="bg-gray-900 border border-gray-800 rounded-lg px-4 py-3">
            <p className="text-xs text-gray-500 mb-1">Vs network avg wait</p>
            <p className="text-2xl font-semibold text-white">
              {formatHeadwayMinutes(routeWait)} min
            </p>
            <p
              className={`text-xs mt-1 ${
                waitDelta.better === true
                  ? 'text-green-400'
                  : waitDelta.better === false
                    ? 'text-red-400'
                    : 'text-gray-500'
              }`}
            >
              {waitDelta.text}
              {' · '}
              network {formatHeadwayMinutes(networkWait)} min
            </p>
          </div>
        )}
      </div>
    </div>
  )
}

export default HeadwayContextPanel
