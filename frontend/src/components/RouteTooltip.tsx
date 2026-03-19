import type { RouteProperties } from '../types/api'

interface Props {
  properties: RouteProperties
  x: number
  y: number
}

const OFFSET = 14

const RouteTooltip = ({ properties, x, y }: Props) => {
  const { routeId, routeName, avgRides } = properties

  return (
    <div
      className="absolute z-10 pointer-events-none"
      style={{ left: x + OFFSET, top: y + OFFSET }}
    >
      <div className="bg-gray-900 border border-gray-700/60 rounded-xl shadow-2xl p-4 w-52">
        <div className="flex items-center gap-2 mb-3">
          <span className="shrink-0 bg-blue-600 text-white text-xs font-mono font-bold px-2 py-0.5 rounded-full">
            {routeId}
          </span>
          <span className="text-white text-sm font-semibold leading-tight truncate">
            {routeName}
          </span>
        </div>

        <div className="h-px bg-gray-700/60 mb-3" />

        <div className="text-gray-400 text-xs uppercase tracking-widest mb-1">
          Avg Daily Ridership
        </div>
        {avgRides != null ? (
          <div className="text-white text-2xl font-bold tabular-nums">
            {Math.round(avgRides).toLocaleString()}
          </div>
        ) : (
          <div className="text-gray-500 text-sm italic">No data available</div>
        )}
      </div>
    </div>
  )
}

export default RouteTooltip
