import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import type { RecoveryBucket } from '../lib/ridershipUtils'

interface Props {
  data: RecoveryBucket[]
  height?: number
}

const RecoveryDistributionChart = ({ data, height = 180 }: Props) => {
  if (!data.some(d => d.count > 0)) {
    return (
      <div className="flex items-center justify-center text-gray-500 text-sm" style={{ height }}>
        No recovery data available
      </div>
    )
  }

  return (
    <ResponsiveContainer width="100%" height={height}>
      <BarChart data={data} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" vertical={false} />
        <XAxis
          dataKey="label"
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
        />
        <Tooltip
          contentStyle={{ background: '#111827', border: '1px solid #374151', borderRadius: 8, fontSize: 12 }}
          labelStyle={{ color: '#9ca3af' }}
          formatter={(value: number) => [`${value} routes`, 'Count']}
        />
        <Bar dataKey="count" fill="#dc2626" radius={[4, 4, 0, 0]} />
      </BarChart>
    </ResponsiveContainer>
  )
}

export default RecoveryDistributionChart
