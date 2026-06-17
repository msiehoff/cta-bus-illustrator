import { useMemo } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  ReferenceArea,
  type TooltipProps,
} from 'recharts'
import type { RidershipDataPoint, RidershipType } from '../types/api'
import { PANDEMIC_PERIOD } from '../lib/ridershipUtils'

export type WindowKey = 'all' | '5y' | '1y'

export const WINDOWS: { key: WindowKey; label: string }[] = [
  { key: '1y',  label: 'Last year'    },
  { key: '5y',  label: 'Last 5 years' },
  { key: 'all', label: 'All time'     },
]

interface Props {
  records: RidershipDataPoint[]
  window: WindowKey
  onWindowChange: (w: WindowKey) => void
  height?: number
  overlayRecords?: RidershipDataPoint[]
  overlayLabel?: string
  highlightType?: RidershipType
  showPandemicBand?: boolean
}

interface ChartRow {
  month: string
  weekday?: number
  saturday?: number
  sunday?: number
  system?: number
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

// Returns the cutoff month string relative to a given anchor month (the latest
// month in the dataset), NOT relative to today. This ensures the window label
// matches what's actually shown — the data may lag the current date by months.
export function cutoffMonth(w: WindowKey, anchorMonth: string): string | null {
  if (w === 'all') return null
  const [year, mon] = anchorMonth.split('-').map(Number)
  const years = w === '1y' ? 1 : 5
  const cutoffYear = year - years
  return `${cutoffYear}-${String(mon).padStart(2, '0')}`
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

export default function RidershipChart({
  records,
  window: w,
  onWindowChange,
  height = 220,
  overlayRecords,
  overlayLabel = 'System avg',
  highlightType,
  showPandemicBand = true,
}: Props) {
  const allData = useMemo<ChartRow[]>(() => {
    const byMonth = new Map<string, ChartRow>()
    for (const r of records) {
      if (!byMonth.has(r.month)) byMonth.set(r.month, { month: r.month })
      const row = byMonth.get(r.month)!
      row[r.type] = Math.round(r.avgRides)
    }

    if (overlayRecords?.length && highlightType) {
      for (const r of overlayRecords) {
        if (r.type !== highlightType) continue
        if (!byMonth.has(r.month)) byMonth.set(r.month, { month: r.month })
        const row = byMonth.get(r.month)! 
        row.system = Math.round(r.avgRides)
      }
    }

    return Array.from(byMonth.values()).sort((a, b) => a.month.localeCompare(b.month))
  }, [records, overlayRecords, highlightType])

  // Anchor to the latest month in the data, not today
  const latestMonth = allData.length ? allData[allData.length - 1].month : ''

  const data = useMemo<ChartRow[]>(() => {
    if (!latestMonth) return allData
    const cutoff = cutoffMonth(w, latestMonth)
    if (!cutoff) return allData
    return allData.filter(row => row.month >= cutoff)
  }, [allData, w, latestMonth])

  const visibleLines = highlightType
    ? LINE_CONFIG.filter(cfg => cfg.key === highlightType && records.some(r => r.type === cfg.key))
    : LINE_CONFIG.filter(cfg => records.some(r => r.type === cfg.key))

  const hasOverlay = Boolean(highlightType && overlayRecords?.some(r => r.type === highlightType))

  const pandemicBandVisible = useMemo(() => {
    if (!showPandemicBand || !data.length) return false
    const first = data[0].month
    const last = data[data.length - 1].month
    return last >= PANDEMIC_PERIOD.start && first <= PANDEMIC_PERIOD.end
  }, [data, showPandemicBand])

  if (allData.length === 0) {
    return (
      <div className="flex items-center justify-center text-gray-500 text-sm" style={{ height }}>
        No data available
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <div className="flex gap-4 flex-wrap">
          {visibleLines.map(cfg => (
            <div key={cfg.key} className="flex items-center gap-1.5">
              <div className="w-2.5 h-2.5 rounded-full" style={{ background: cfg.color }} />
              <span className="text-xs text-gray-400 capitalize">{cfg.label}</span>
            </div>
          ))}
          {hasOverlay && (
            <div className="flex items-center gap-1.5">
              <div className="w-2.5 h-0.5 bg-gray-400" />
              <span className="text-xs text-gray-400">{overlayLabel}</span>
            </div>
          )}
          {pandemicBandVisible && (
            <div className="flex items-center gap-1.5">
              <div className="w-2.5 h-2.5 rounded-sm bg-gray-700/80" />
              <span className="text-xs text-gray-400">Pandemic period</span>
            </div>
          )}
        </div>
        <div className="flex items-center gap-1">
          {WINDOWS.map(win => (
            <button
              key={win.key}
              onClick={() => onWindowChange(win.key)}
              className={`px-2.5 py-1 rounded-md text-xs transition-colors ${
                w === win.key
                  ? 'bg-gray-700 text-white'
                  : 'text-gray-500 hover:text-gray-300 hover:bg-gray-800'
              }`}
            >
              {win.label}
            </button>
          ))}
        </div>
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
          {pandemicBandVisible && (
            <ReferenceArea
              x1={PANDEMIC_PERIOD.start}
              x2={PANDEMIC_PERIOD.end}
              strokeOpacity={0}
              fill="#374151"
              fillOpacity={0.35}
            />
          )}
          {visibleLines.map(cfg => (
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
          {hasOverlay && (
            <Line
              type="monotone"
              dataKey="system"
              stroke="#9ca3af"
              strokeWidth={2}
              strokeDasharray="6 4"
              dot={false}
              activeDot={{ r: 4, strokeWidth: 0 }}
            />
          )}
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
