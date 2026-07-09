import { useEffect, useState } from 'react'
import { twMerge } from 'tailwind-merge'
import type { RidershipType } from '../types/api'
import type { RankedEntry, TooltipData } from './RouteMap'
import { legendGradient, LEGEND_TICKS } from '../lib/ridershipColors'
import { formatMonth } from '../lib/ridershipUtils'
import useMediaQuery from '../hooks/useMediaQuery'
import { CloseIcon } from './SVG'
import MonthSelector from './MonthSelector'

interface Props {
  availableMonths: string[]
  selectedMonth: string | null
  ridershipType: RidershipType
  rankedEntries: RankedEntry[]
  onMonthChange: (month: string) => void
  onTypeChange: (type: RidershipType) => void
  onRouteClick: (entry: RankedEntry) => void
  onExpandedChange?: (expanded: boolean) => void
}

const RIDERSHIP_TYPES: { value: RidershipType; label: string }[] = [
  { value: 'weekday', label: 'Weekday' },
  { value: 'saturday', label: 'Saturday' },
  { value: 'sunday', label: 'Sunday' },
]

const formatRidesCompact = (n: number): string => {
  if (n >= 1000) return `${(n / 1000).toFixed(1).replace(/\.0$/, '')}k`
  return String(Math.round(n))
}

const getEntryDisplay = (data: TooltipData) =>
  data.type === 'single'
    ? { routeId: data.properties.routeId, name: data.properties.routeName }
    : { routeId: data.local.routeId, name: data.local.routeName }

const ridershipTypeLabel = (type: RidershipType) =>
  RIDERSHIP_TYPES.find(t => t.value === type)?.label ?? type

const selectClass = 'w-full bg-gray-800 border border-gray-700 text-white text-sm rounded-lg px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:opacity-50 cursor-pointer'
const labelClass = 'text-gray-400 text-xs uppercase tracking-widest block mb-1.5'
const dividerClass = 'h-px bg-gray-700/60'
const LEGEND_TIP_DISMISSED_KEY = 'cta.legendTipDismissed'

const FilterBar = ({
  availableMonths,
  selectedMonth,
  ridershipType,
  rankedEntries,
  onMonthChange,
  onTypeChange,
  onRouteClick,
  onExpandedChange,
}: Props) => {
  const isMobile = useMediaQuery('(max-width: 767px)')
  const [expanded, setExpanded] = useState(false)
  const [showLegendTip, setShowLegendTip] = useState(false)

  useEffect(() => {
    setExpanded(!isMobile)
  }, [isMobile])

  useEffect(() => {
    onExpandedChange?.(expanded)
  }, [expanded, onExpandedChange])

  useEffect(() => {
    const isDismissed = localStorage.getItem(LEGEND_TIP_DISMISSED_KEY) === '1'
    setShowLegendTip(!isDismissed)
  }, [])

  const dismissLegendTip = () => {
    localStorage.setItem(LEGEND_TIP_DISMISSED_KEY, '1')
    setShowLegendTip(false)
  }

  const summaryLabel = selectedMonth
    ? `${ridershipTypeLabel(ridershipType)} · ${formatMonth(selectedMonth)}`
    : 'Filters'

  if (isMobile && !expanded) {
    return (
      <div className="absolute top-3 left-3 z-10 pointer-events-auto">
        <button
          type="button"
          onClick={() => setExpanded(true)}
          className="bg-gray-900/90 backdrop-blur border border-gray-700/60 rounded-xl shadow-2xl px-3 py-2 text-left"
        >
          <span className="text-[10px] uppercase tracking-widest text-gray-500 block">Filters</span>
          <span className="text-sm text-white font-medium">{summaryLabel}</span>
        </button>
      </div>
    )
  }

  return (
    <div className="absolute top-3 left-3 z-10 pointer-events-auto">
      <div className={twMerge(
        'bg-gray-900/90 backdrop-blur border border-gray-700/60 rounded-xl shadow-2xl p-3 flex flex-col gap-3 w-56',
        isMobile && 'max-h-[calc(100vh-5rem)] overflow-y-auto',
      )}>
        {isMobile && (
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-gray-200">Filters</span>
            <button
              type="button"
              onClick={() => setExpanded(false)}
              aria-label="Close filters"
              className="flex items-center justify-center w-8 h-8 rounded-lg text-gray-400 hover:text-white hover:bg-gray-800 transition-colors"
            >
              <CloseIcon size={18} />
            </button>
          </div>
        )}

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

        <MonthSelector
          months={availableMonths}
          selectedMonth={selectedMonth}
          onMonthChange={onMonthChange}
          disabled={availableMonths.length === 0}
          compact
        />

        <div className={dividerClass} />

        <div>
          <span className="text-gray-200 text-xs font-semibold block">Average daily riders</span>
          {selectedMonth ? (
            <span className="text-gray-500 text-[10px] block mt-0.5 mb-1.5 tabular-nums">
              {ridershipTypeLabel(ridershipType)} · {formatMonth(selectedMonth)}
            </span>
          ) : (
            <span className="text-gray-500 text-[10px] block mt-0.5 mb-1.5">
              Select month and day type above
            </span>
          )}
          <p className="text-white text-[11px] leading-snug mb-2">
            Thicker lines and warmer colors mean more riders on this scale.
          </p>
          {showLegendTip && (
            <div className="mb-2 rounded-md border border-blue-400/40 bg-blue-500/10 px-2 py-1.5">
              <p className="text-[10px] text-blue-100 leading-snug">
                {isMobile
                  ? 'Tip: tap any line to see route details. Tap again to open the route page.'
                  : 'Tip: hover any line to see route details and exact ridership.'}
              </p>
              <button
                type="button"
                onClick={dismissLegendTip}
                className="mt-1 text-[10px] text-blue-200 hover:text-white transition-colors"
              >
                Got it
              </button>
            </div>
          )}
          <div className="h-2.5 rounded-full w-full" style={{ background: legendGradient }} />
          <div className="flex justify-between mt-1">
            {LEGEND_TICKS.map((tick, i) => (
              <span key={i} className="text-gray-400 text-xs tabular-nums">{tick}</span>
            ))}
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
}

export default FilterBar
