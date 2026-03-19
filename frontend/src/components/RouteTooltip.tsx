import type { TooltipData } from './RouteMap'
import type { RouteProperties } from '../types/api'

interface Props {
  data: TooltipData
  x: number
  y: number
  rank?: number
}

const OFFSET = 14

const RouteIdBadge = ({ routeId }: { routeId: string }) => (
  <span className="shrink-0 bg-blue-600 text-white text-xs font-mono font-bold px-2 py-0.5 rounded-full">
    {routeId}
  </span>
)

const Divider = () => <div className="h-px bg-gray-700/60 my-3" />

const RankBanner = ({ rank }: { rank: number }) => (
  <div className="flex items-center gap-1.5 bg-amber-500/15 border border-amber-500/30 rounded-lg px-2.5 py-1.5 mb-3">
    <span className="text-amber-400 text-xs font-bold tabular-nums">#{rank}</span>
    <span className="text-amber-300/80 text-xs">highest avg ridership</span>
  </div>
)

const RidershipRow = ({ properties }: { properties: RouteProperties }) => (
  <div className="flex items-center justify-between gap-2 py-1">
    <div className="flex items-center gap-2 min-w-0">
      <RouteIdBadge routeId={properties.routeId} />
      <span className="text-gray-300 text-xs truncate">{properties.routeName}</span>
    </div>
    <span className="text-white text-sm font-semibold tabular-nums shrink-0">
      {properties.avgRides != null ? Math.round(properties.avgRides).toLocaleString() : '—'}
    </span>
  </div>
)

const SingleView = ({ properties }: { properties: RouteProperties }) => {
  const { routeId, routeName, avgRides } = properties
  return (
    <>
      <div className="flex items-center gap-2 mb-3">
        <RouteIdBadge routeId={routeId} />
        <span className="text-white text-sm font-semibold leading-tight truncate">{routeName}</span>
      </div>
      <Divider />
      <div className="text-gray-400 text-xs uppercase tracking-widest mb-1">Avg Daily Ridership</div>
      {avgRides != null ? (
        <div className="text-white text-2xl font-bold tabular-nums">{Math.round(avgRides).toLocaleString()}</div>
      ) : (
        <div className="text-gray-500 text-sm italic">No data available</div>
      )}
    </>
  )
}

const CorridorView = ({ local, express }: { local: RouteProperties; express: RouteProperties }) => {
  const localRides = local.avgRides ?? 0
  const expressRides = express.avgRides ?? 0
  const total = localRides + expressRides
  const corridorName = local.routeName.replace(/ Express$/i, '')

  return (
    <>
      <div className="text-white text-sm font-semibold leading-tight mb-1">{corridorName} Corridor</div>
      <div className="text-gray-400 text-xs">Express + Local</div>
      <Divider />
      <div className="text-gray-400 text-xs uppercase tracking-widest mb-1">Total Avg Daily Ridership</div>
      <div className="text-white text-2xl font-bold tabular-nums mb-3">
        {total > 0 ? total.toLocaleString() : '—'}
      </div>
      <Divider />
      <div className="text-gray-400 text-xs uppercase tracking-widest mb-1">Breakdown</div>
      <RidershipRow properties={local} />
      <RidershipRow properties={express} />
    </>
  )
}

const RouteTooltip = ({ data, x, y, rank }: Props) => (
  <div
    className="absolute z-10 pointer-events-none"
    style={{ left: x + OFFSET, top: y + OFFSET }}
  >
    <div className="bg-gray-900 border border-gray-700/60 rounded-xl shadow-2xl p-4 w-60">
      {rank !== undefined && <RankBanner rank={rank} />}
      {data.type === 'single'
        ? <SingleView properties={data.properties} />
        : <CorridorView local={data.local} express={data.express} />
      }
    </div>
  </div>
)

export default RouteTooltip
