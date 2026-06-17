import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import type { WeekendSharePoint } from '../lib/ridershipUtils'

interface Props {
  data: WeekendSharePoint[]
  height?: number
}

const WeekendShareChart = ({ data, height = 180 }: Props) => {
  if (!data.length) {
    return (
      <div className="flex items-center justify-center text-gray-500 text-sm" style={{ height }}>
        No weekend share data available
      </div>
    )
  }

  return (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={data} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
        <XAxis
          dataKey="month"
          tick={{ fill: '#9ca3af', fontSize: 11 }}
          tickLine={false}
          axisLine={false}
          interval="preserveStartEnd"
        />
        <YAxis
          tickFormatter={v => `${v.toFixed(0)}%`}
          tick={{ fill: '#9ca3af', fontSize: 11 }}
          tickLine={false}
          axisLine={false}
          width={44}
        />
        <Tooltip
          contentStyle={{ background: '#111827', border: '1px solid #374151', borderRadius: 8, fontSize: 12 }}
          labelStyle={{ color: '#9ca3af' }}
          formatter={(value: number) => [`${value.toFixed(1)}%`, 'Weekend / weekday']}
        />
        <Line
          type="monotone"
          dataKey="sharePct"
          stroke="#f59e0b"
          strokeWidth={2}
          dot={false}
          activeDot={{ r: 4, strokeWidth: 0 }}
        />
      </LineChart>
    </ResponsiveContainer>
  )
}

export default WeekendShareChart
