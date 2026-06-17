import { useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  CartesianGrid,
  ResponsiveContainer,
  Scatter,
  ScatterChart,
  Tooltip,
  XAxis,
  YAxis,
  ZAxis,
} from 'recharts'
import type { ScatterPoint } from '../lib/ridershipUtils'
import { formatRides } from '../lib/ridershipUtils'

interface Props {
  data: ScatterPoint[]
  height?: number
}

const RecoveryScatterChart = ({ data, height = 220 }: Props) => {
  const navigate = useNavigate()

  const yDomain = useMemo(() => {
    if (!data.length) return [0, 120]
    const min = Math.min(...data.map(d => d.recoveryPct))
    const max = Math.max(...data.map(d => d.recoveryPct))
    return [Math.max(0, Math.floor(min / 10) * 10 - 10), Math.ceil(max / 10) * 10 + 10]
  }, [data])

  if (!data.length) {
    return (
      <div className="flex items-center justify-center text-gray-500 text-sm" style={{ height }}>
        No recovery data available
      </div>
    )
  }

  return (
    <ResponsiveContainer width="100%" height={height}>
      <ScatterChart margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
        <XAxis
          type="number"
          dataKey="current"
          name="Ridership"
          tickFormatter={formatRides}
          tick={{ fill: '#9ca3af', fontSize: 11 }}
          tickLine={false}
          axisLine={false}
        />
        <YAxis
          type="number"
          dataKey="recoveryPct"
          name="Recovery"
          domain={yDomain}
          tickFormatter={v => `${v}%`}
          tick={{ fill: '#9ca3af', fontSize: 11 }}
          tickLine={false}
          axisLine={false}
          width={44}
        />
        <ZAxis type="number" dataKey="current" range={[30, 120]} />
        <Tooltip
          cursor={{ strokeDasharray: '3 3' }}
          content={({ active, payload }) => {
            if (!active || !payload?.length) return null
            const point = payload[0].payload as ScatterPoint
            return (
              <div className="bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-xs shadow-lg">
                <p className="text-white font-medium">{point.routeName}</p>
                <p className="text-gray-400 mt-1">{formatRides(point.current)} riders</p>
                <p className="text-gray-400">{point.recoveryPct.toFixed(0)}% of pre-COVID</p>
              </div>
            )
          }}
        />
        <Scatter
          data={data}
          fill="#dc2626"
          fillOpacity={0.75}
          onClick={point => navigate(`/routes/${point.routeId}`, { state: { routeName: point.routeName } })}
          className="cursor-pointer"
        />
      </ScatterChart>
    </ResponsiveContainer>
  )
}

export default RecoveryScatterChart
