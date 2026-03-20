import { twMerge } from 'tailwind-merge'
import type { RidershipType } from '../types/api'
import type { RankedEntry, TooltipData } from './RouteMap'
import { legendGradient, LEGEND_TICKS, NO_DATA_COLOR } from '../lib/ridershipColors'

interface Props {
  availableMonths: string[]
  selectedMonth: string | null
  ridershipType: RidershipType
  rankedEntries: RankedEntry[]
  onMonthChange: (month: string) => void
  onTypeChange: (type: RidershipType) => void
  onRouteClick: (entry: RankedEntry) => void
}

const RIDERSHIP_TYPES: { value: RidershipType; label: string }[] = [
  { value: 'weekday', label: 'Weekday' },
  { value: 'saturday', label: 'Saturday' },
  { value: 'sunday', label: 'Sunday' },
]

const formatMonth = (yyyyMM: string) => {
  const [year, month] = yyyyMM.split('-').map(Number)
  return new Date(year, month - 1, 1).toLocaleString('default', { month: 'long', year: 'numeric' })
}

const formatRidesCompact = (n: number): string => {
  if (n >= 1000) return `${(n / 1000).toFixed(1).replace(/\.0$/, '')}k`
  return String(Math.round(n))
}

const getEntryDisplay = (data: TooltipData) =>
  data.type === 'single'
    ? { routeId: data.properties.routeId, name: data.properties.routeName }
    : { routeId: data.local.routeId, name: data.local.routeName }

const selectClass = 'w-full bg-gray-800 border border-gray-700 text-white text-sm rounded-lg px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:opacity-50 cursor-pointer'
const labelClass = 'text-gray-400 text-xs uppercase tracking-widest block mb-1.5'
const dividerClass = 'h-px bg-gray-700/60'

const FilterBar = ({ availableMonths, selectedMonth, ridershipType, rankedEntries, onMonthChange, onTypeChange, onRouteClick }: Props) => (
  <div className="absolute top-3 left-3 z-10 pointer-events-auto">
    <div className="bg-gray-900/90 backdrop-blur border border-gray-700/60 rounded-xl shadow-2xl p-3 flex flex-col gap-3 w-48">

      <div>
        <label className={labelClass}>Ridership</label>
        <select
          value={ridershipType}
          onChange={e => onTypeChange(e.target.value as RidershipType)}
          className={selectClass}
        >
          {RIDERSHIP_TYPES.map(({ value, label }) => (
            <option key={value} value={value}>{label}</option>
          ))}
        </select>
      </div>

      <div>
        <label className={labelClass}>Month</label>
        <select
          value={selectedMonth ?? ''}
          onChange={e => onMonthChange(e.target.value)}
          disabled={availableMonths.length === 0}
          className={selectClass}
        >
          {availableMonths.length === 0
            ? <option value="">No data</option>
            : availableMonths.map(m => <option key={m} value={m}>{formatMonth(m)}</option>)
          }
        </select>
      </div>

      <div className={dividerClass} />

      <div>
        <label className={labelClass}>Avg Daily Riders</label>
        <div className="h-2.5 rounded-full w-full" style={{ background: legendGradient }} />
        <div className="flex justify-between mt-1">
          {LEGEND_TICKS.map((tick, i) => (
            <span key={i} className="text-gray-400 text-xs tabular-nums">{tick}</span>
          ))}
        </div>
        <div className="flex items-center gap-1.5 mt-2">
          <div className="w-3 h-2.5 rounded-sm shrink-0" style={{ backgroundColor: NO_DATA_COLOR }} />
          <span className="text-gray-500 text-xs">No data</span>
        </div>
      </div>

      {rankedEntries.length > 0 && (
        <>
          <div className={dividerClass} />

          <div>
            <label className={labelClass}>Top Routes</label>
            <div className="flex flex-col gap-0.5">
              {rankedEntries.map(entry => {
                const { routeId, name } = getEntryDisplay(entry.data)
                return (
                  <div
                    key={entry.rank}
                    className="flex items-center gap-1.5 py-1 rounded-lg px-1 cursor-pointer hover:bg-gray-700/50 transition-colors"
                    onClick={() => onRouteClick(entry)}
                  >
                    <span className="text-gray-500 text-xs tabular-nums w-4 shrink-0 text-right">
                      {entry.rank}
                    </span>
                    <span className={twMerge(
                      'shrink-0 inline-flex items-center bg-blue-600 text-white font-mono font-bold rounded-full px-1.5 py-0.5 leading-none',
                      'text-[10px]'
                    )}>
                      {routeId}
                    </span>
                    <span className="text-gray-300 text-xs truncate flex-1 min-w-0">{name}</span>
                    <span className="text-white text-xs tabular-nums shrink-0 font-medium">
                      {formatRidesCompact(entry.totalRides)}
                    </span>
                  </div>
                )
              })}
            </div>
          </div>
        </>
      )}

    </div>
  </div>
)

export default FilterBar
