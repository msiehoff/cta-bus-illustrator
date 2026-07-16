import { useMemo, useState } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  type TooltipProps,
} from 'recharts'
import type { HeadwayDayPoint } from '../types/api'
import { formatHeadwayMinutes } from '../lib/headwayUtils'

type MetricKey = 'medianMinutes' | 'avgWaitMinutes'

interface Props {
  series: HeadwayDayPoint[]
  height?: number
}

const METRICS: { key: MetricKey; label: string; color: string }[] = [
  { key: 'medianMinutes', label: 'Median headway', color: '#dc2626' },
  { key: 'avgWaitMinutes', label: 'Avg wait', color: '#f97316' },
]

const CustomTooltip = ({ active, payload, label }: TooltipProps<number, string>) => {
  if (!active || !payload?.length) return null
  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-xs shadow-lg">
      <p className="text-gray-400 mb-1.5">{label}</p>
      {payload.map(entry => (
        <div key={String(entry.dataKey)} className="flex items-center gap-2 mb-0.5">
          <div className="w-2 h-2 rounded-full" style={{ background: entry.color }} />
          <span className="text-gray-300">{entry.name}:</span>
          <span className="text-white font-medium">
            {formatHeadwayMinutes(entry.value ?? 0)} min
          </span>
        </div>
      ))}
    </div>
  )
}

const HeadwayDailyChart = ({ series, height = 220 }: Props) => {
  const [metric, setMetric] = useState<MetricKey>('medianMinutes')
  const active = METRICS.find(m => m.key === metric) ?? METRICS[0]

  const data = useMemo(
    () =>
      series.map(p => ({
        date: p.serviceDate,
        medianMinutes: p.medianMinutes,
        avgWaitMinutes: p.avgWaitMinutes,
      })),
    [series],
  )

  if (!series.length) {
    return (
      <div className="text-gray-500 text-sm py-10 text-center">
        No daily points in this period yet.
      </div>
    )
  }

  return (
    <div>
      <div className="flex gap-2 mb-3">
        {METRICS.map(m => (
          <button
            key={m.key}
            type="button"
            onClick={() => setMetric(m.key)}
            className={`text-xs px-2.5 py-1 rounded border transition-colors ${
              metric === m.key
                ? 'border-red-800 bg-red-950/50 text-red-300'
                : 'border-gray-800 text-gray-500 hover:text-gray-300'
            }`}
          >
            {m.label}
          </button>
        ))}
      </div>
      <ResponsiveContainer width="100%" height={height}>
        <LineChart data={data} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
          <CartesianGrid stroke="#1f2937" strokeDasharray="3 3" />
          <XAxis
            dataKey="date"
            tick={{ fill: '#6b7280', fontSize: 11 }}
            tickLine={false}
            axisLine={{ stroke: '#374151' }}
          />
          <YAxis
            tick={{ fill: '#6b7280', fontSize: 11 }}
            tickLine={false}
            axisLine={false}
            tickFormatter={v => `${v}`}
            width={36}
            unit="m"
          />
          <Tooltip content={<CustomTooltip />} />
          <Line
            type="monotone"
            dataKey={active.key}
            name={active.label}
            stroke={active.color}
            strokeWidth={2}
            dot={{ r: 3, fill: active.color }}
            activeDot={{ r: 5 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}

export default HeadwayDailyChart
