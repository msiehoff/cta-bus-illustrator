import { useMemo } from 'react'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import {
  buildRidershipDensity,
  formatRides,
  getRidershipPercentile,
} from '../lib/ridershipUtils'

interface RoutePoint {
  routeId: string
  routeName: string
  current: number
}

interface Props {
  routes: RoutePoint[]
  highlightRouteId?: string
  highlightRouteName?: string
  ridershipType?: string
  height?: number
}

const RidershipDistributionChart = ({
  routes,
  highlightRouteId,
  highlightRouteName,
  ridershipType = 'weekday',
  height = 200,
}: Props) => {
  const bins = useMemo(() => buildRidershipDensity(routes), [routes])

  const highlight = useMemo(() => {
    if (!highlightRouteId) return null
    const route = routes.find(r => r.routeId === highlightRouteId)
    if (!route) return null
    return {
      ridership: route.current,
      percentile: getRidershipPercentile(routes, route.current),
      name: highlightRouteName ?? route.routeName,
    }
  }, [routes, highlightRouteId, highlightRouteName])

  if (!bins.length) {
    return (
      <div className="flex items-center justify-center text-gray-500 text-sm" style={{ height }}>
        Not enough route data
      </div>
    )
  }

  return (
    <div>
      {highlight && (
        <p className="text-sm text-gray-300 mb-3">
          <span className="text-amber-400 font-medium">{highlight.name}</span>
          {' '}carries more {ridershipType} riders than{' '}
          <span className="text-white font-medium">{highlight.percentile}%</span> of routes
          ({formatRides(highlight.ridership)} avg daily)
        </p>
      )}
      <p className="text-xs text-gray-500 mb-3">
        {highlight
          ? 'Height shows how many routes fall in each ridership band — peaks mean lots of routes at that level.'
          : `How avg daily ${ridershipType} ridership is spread across routes. Height = number of routes in each band.`}
      </p>
      <ResponsiveContainer width="100%" height={height}>
        <AreaChart data={bins} margin={{ top: 8, right: 12, left: 4, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" vertical={false} />
          <XAxis
            dataKey="ridership"
            type="number"
            domain={['dataMin', 'dataMax']}
            tickFormatter={formatRides}
            tick={{ fill: '#9ca3af', fontSize: 11 }}
            tickLine={false}
            axisLine={false}
          />
          <YAxis
            allowDecimals={false}
            tick={{ fill: '#9ca3af', fontSize: 11 }}
            tickLine={false}
            axisLine={false}
            width={28}
            label={{
              value: 'Routes',
              angle: -90,
              position: 'insideLeft',
              fill: '#6b7280',
              fontSize: 10,
              offset: 10,
            }}
          />
          <Tooltip
            content={({ active, payload }) => {
              if (!active || !payload?.length) return null
              const point = payload[0].payload as { ridership: number; routeCount: number; rangeLabel: string }
              return (
                <div className="bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-xs shadow-lg">
                  <p className="text-gray-400">{point.rangeLabel} avg daily</p>
                  <p className="text-white font-medium">
                    {point.routeCount} route{point.routeCount === 1 ? '' : 's'}
                  </p>
                </div>
              )
            }}
          />
          <Area
            type="monotone"
            dataKey="routeCount"
            stroke="#dc2626"
            strokeWidth={2}
            fill="#dc2626"
            fillOpacity={0.2}
            isAnimationActive={false}
          />
          {highlight && (
            <ReferenceLine
              x={highlight.ridership}
              stroke="#fbbf24"
              strokeWidth={2}
              strokeDasharray="4 4"
              label={{
                value: highlight.name,
                position: 'top',
                fill: '#fbbf24',
                fontSize: 11,
              }}
            />
          )}
        </AreaChart>
      </ResponsiveContainer>
      {highlight && (
        <div className="flex items-center gap-1.5 mt-2">
          <span className="inline-block w-4 border-t-2 border-dashed border-amber-400" />
          <span className="text-[10px] text-gray-500">This route</span>
        </div>
      )}
    </div>
  )
}

export default RidershipDistributionChart
