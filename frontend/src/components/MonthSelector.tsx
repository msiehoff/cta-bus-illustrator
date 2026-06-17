import { twMerge } from 'tailwind-merge'
import { formatMonth } from '../lib/ridershipUtils'
import {
  formatMonthOption,
  getAdjacentMonth,
  getAvailableMonthsInYear,
  getAvailableYears,
  parseMonth,
  resolveMonthForYear,
} from '../lib/monthSelection'

interface Props {
  months: string[]
  selectedMonth: string | null
  onMonthChange: (month: string) => void
  disabled?: boolean
  compact?: boolean
}

const selectClass =
  'bg-gray-800 border border-gray-700 text-white text-sm rounded-lg px-3 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:opacity-50 cursor-pointer'

const labelClass = 'text-[10px] uppercase tracking-widest text-gray-500 block mb-1.5'

const stepperBtnClass =
  'flex items-center justify-center w-8 h-[34px] rounded-lg border border-gray-700 bg-gray-800 text-gray-400 hover:text-white hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors'

const MonthSelector = ({
  months,
  selectedMonth,
  onMonthChange,
  disabled = false,
  compact = false,
}: Props) => {
  const parsed = parseMonth(selectedMonth)
  const years = getAvailableYears(months)
  const selectedYear = parsed?.year ?? years[0]
  const monthsInYear = selectedYear ? getAvailableMonthsInYear(months, selectedYear) : []

  const canGoNewer = Boolean(selectedMonth && getAdjacentMonth(months, selectedMonth, 'newer'))
  const canGoOlder = Boolean(selectedMonth && getAdjacentMonth(months, selectedMonth, 'older'))

  const handleYearChange = (year: number) => {
    const next = resolveMonthForYear(months, year, parsed?.monthNum)
    if (next) onMonthChange(next)
  }

  return (
    <div className={twMerge('flex flex-wrap items-end gap-2', compact && 'flex-col items-stretch gap-3')}>
      <div className={twMerge('flex items-end gap-2', compact && 'w-full')}>
        <div className={compact ? 'flex-1' : undefined}>
          <label className={labelClass}>Year</label>
          <select
            value={selectedYear ?? ''}
            onChange={e => handleYearChange(Number(e.target.value))}
            disabled={disabled || !years.length}
            className={twMerge(selectClass, compact && 'w-full')}
          >
            {years.map(year => (
              <option key={year} value={year}>{year}</option>
            ))}
          </select>
        </div>
        <div className={compact ? 'flex-1' : undefined}>
          <label className={labelClass}>Month</label>
          <select
            value={selectedMonth ?? ''}
            onChange={e => onMonthChange(e.target.value)}
            disabled={disabled || !monthsInYear.length}
            className={twMerge(selectClass, compact && 'w-full')}
          >
            {monthsInYear.map(month => (
              <option key={month} value={month}>{formatMonthOption(month)}</option>
            ))}
          </select>
        </div>
      </div>

      <div className={twMerge('flex items-center gap-1', compact && 'w-full justify-between')}>
        <button
          type="button"
          aria-label="Previous month"
          disabled={disabled || !canGoOlder}
          onClick={() => {
            const prev = getAdjacentMonth(months, selectedMonth, 'older')
            if (prev) onMonthChange(prev)
          }}
          className={stepperBtnClass}
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <polyline points="15 18 9 12 15 6" />
          </svg>
        </button>
        <span className={twMerge(
          'text-xs text-gray-400 px-2 min-w-[5.5rem] text-center',
          compact && 'flex-1',
        )}>
          {selectedMonth ? formatMonth(selectedMonth) : '—'}
        </span>
        <button
          type="button"
          aria-label="Next month"
          disabled={disabled || !canGoNewer}
          onClick={() => {
            const next = getAdjacentMonth(months, selectedMonth, 'newer')
            if (next) onMonthChange(next)
          }}
          className={stepperBtnClass}
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>
      </div>
    </div>
  )
}

export default MonthSelector
