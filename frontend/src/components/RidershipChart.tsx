import { useMemo } from 'react'
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
import type { RidershipDataPoint } from '../types/api'

interface Props {
  records: RidershipDataPoint[]
  height?: number
}

// One row per month in the chart, with weekday/saturday/sunday as columns.
interface ChartRow {
  month: string
  weekday?: number
  saturday?: number
  sunday?: number
}

const LINE_CONFIG = [
  { key: 'weekday',  color: '#dc2626', label: 'Weekday'  },
  { key: 'saturday', color: '#f97316', label: 'Saturday' },
  { key: 'sunday',   color: '#f59e0b', label: 'Sunday'   },
] as const

function formatRides(value: number): string {
  if (value >= 1000) return `${(value / 1000).toFixed(1)}k`
  return value.toFixed(0)
}

function CustomTooltip({ active, payload, label }: TooltipProps<number, string>) {
  if (!active || !payload?.length) return null
  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-xs shadow-lg">
      <p className="text-gray-400 mb-1.5">{label}</p>
      {payload.map(entry => (
        <div key={entry.dataKey} className="flex items-center gap-2 mb-0.5">
          <div className="w-2 h-2 rounded-full" style={{ background: entry.color }} />
          <span className="text-gray-300 capitalize">{entry.dataKey}:</span>
          <span className="text-white font-medium">{formatRides(entry.value ?? 0)}</span>
        </div>
      ))}
    </div>
  )
}

export default function RidershipChart({ records, height = 220 }: Props) {
  const data = useMemo<ChartRow[]>(() => {
    const byMonth = new Map<string, ChartRow>()
    for (const r of records) {
      if (!byMonth.has(r.month)) byMonth.set(r.month, { month: r.month })
      const row = byMonth.get(r.month)!
      row[r.type] = Math.round(r.avgRides)
    }
    // Sort chronologically
    return Array.from(byMonth.values()).sort((a, b) => a.month.localeCompare(b.month))
  }, [records])

  // Determine which day types actually have data
  const activeLines = LINE_CONFIG.filter(cfg =>
    records.some(r => r.type === cfg.key)
  )

  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center text-gray-500 text-sm" style={{ height }}>
        No data available
      </div>
    )
  }

  return (
    <div>
      {/* Legend */}
      <div className="flex gap-4 mb-3">
        {activeLines.map(cfg => (
          <div key={cfg.key} className="flex items-center gap-1.5">
            <div className="w-2.5 h-2.5 rounded-full" style={{ background: cfg.color }} />
            <span className="text-xs text-gray-400 capitalize">{cfg.label}</span>
          </div>
        ))}
      </div>

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
            tickFormatter={formatRides}
            tick={{ fill: '#9ca3af', fontSize: 11 }}
            tickLine={false}
            axisLine={false}
            width={44}
          />
          <Tooltip content={<CustomTooltip />} />
          {activeLines.map(cfg => (
            <Line
              key={cfg.key}
              type="monotone"
              dataKey={cfg.key}
              stroke={cfg.color}
              strokeWidth={2}
              dot={false}
              activeDot={{ r: 4, strokeWidth: 0 }}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
