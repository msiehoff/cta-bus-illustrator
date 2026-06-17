import { useState } from 'react'
import {
  formatMonth,
  formatRecoveryPct,
  formatRecoveryPlainEnglish,
  formatRides,
  recoveryBarColorClass,
  recoveryColorClass,
} from '../lib/ridershipUtils'

interface Props {
  title?: string
  currentMonth: string
  benchmarkMonth: string
  current: number
  preCovid?: number
  recovery?: number
  dayTypeLabel?: string
}

export const PrePandemicExplainer = () => {
  const [open, setOpen] = useState(false)

  return (
    <div className="mt-3 pt-3 border-t border-gray-800">
      <button
        type="button"
        onClick={() => setOpen(v => !v)}
        className="text-xs text-gray-500 hover:text-gray-300 transition-colors"
      >
        {open ? 'Hide' : 'Why this benchmark?'}
      </button>
      {open && (
        <ul className="mt-2 space-y-1.5 text-xs text-gray-500 list-disc pl-4">
          <li>
            <span className="text-gray-400">2019</span> is the last normal year before COVID-19
            sharply reduced transit ridership in 2020.
          </li>
          <li>
            We compare the <span className="text-gray-400">same calendar month</span> so seasonal
            patterns (summer vs winter) don&apos;t skew the result.
          </li>
          <li>
            Each <span className="text-gray-400">day type</span> is compared separately because
            weekday commuter routes often recover differently from weekend service.
          </li>
        </ul>
      )}
    </div>
  )
}

const RecoveryBanner = ({
  title = 'Ridership vs pre-pandemic',
  currentMonth,
  benchmarkMonth,
  current,
  preCovid,
  recovery,
  dayTypeLabel = 'weekday',
}: Props) => {
  if (recovery == null || preCovid == null) return null

  const barWidth = Math.min(recovery, 100)

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg px-5 py-4 mb-5">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="min-w-0">
          <p className="text-xs text-gray-500 uppercase tracking-wide">{title}</p>
          <p className="text-2xl font-semibold text-white mt-0.5">
            <span className={recoveryColorClass(recovery)}>
              {formatRecoveryPct(recovery)}
            </span>
            <span className="text-gray-400 text-base font-normal ml-2">
              of pre-pandemic {dayTypeLabel} levels
            </span>
          </p>
          <p className="text-sm text-gray-400 mt-1.5">
            {formatRecoveryPlainEnglish(recovery)}
          </p>
          <p className="text-xs text-gray-500 mt-2">
            Transit ridership dropped sharply in 2020. This shows how close current ridership is
            to the same month in 2019, before the pandemic.
          </p>
          <p className="text-xs text-gray-500 mt-1">
            {formatRides(current)} avg {dayTypeLabel} riders · {formatMonth(currentMonth)}
            {' · '}
            {formatRides(preCovid)} pre-pandemic ({formatMonth(benchmarkMonth)})
          </p>
          <PrePandemicExplainer />
        </div>
        <div className="sm:w-48 shrink-0">
          <div className="h-2 bg-gray-800 rounded-full overflow-hidden">
            <div
              className={`h-full rounded-full transition-all ${recoveryBarColorClass(recovery)}`}
              style={{ width: `${barWidth}%` }}
            />
          </div>
          <p className="text-[10px] text-gray-600 mt-1 text-right">
            100% = same month in 2019, before COVID
          </p>
        </div>
      </div>
    </div>
  )
}

export default RecoveryBanner
