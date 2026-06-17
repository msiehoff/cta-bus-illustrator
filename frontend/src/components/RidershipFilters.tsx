import type { RidershipType } from '../types/api'
import { formatMonth } from '../lib/ridershipUtils'

const RIDERSHIP_TYPES: { value: RidershipType; label: string }[] = [
  { value: 'weekday', label: 'Weekday' },
  { value: 'saturday', label: 'Saturday' },
  { value: 'sunday', label: 'Sunday' },
]

interface Props {
  months: string[]
  selectedMonth: string | null
  ridershipType: RidershipType
  monthsLoading?: boolean
  onMonthChange: (month: string) => void
  onTypeChange: (type: RidershipType) => void
}

const selectClass =
  'bg-gray-800 border border-gray-700 text-white text-sm rounded-lg px-3 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:opacity-50 cursor-pointer'

const RidershipFilters = ({
  months,
  selectedMonth,
  ridershipType,
  monthsLoading = false,
  onMonthChange,
  onTypeChange,
}: Props) => (
  <div className="flex flex-wrap items-end gap-3 mb-5">
    <div>
      <label className="text-[10px] uppercase tracking-widest text-gray-500 block mb-1.5">
        Day type
      </label>
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
      <label className="text-[10px] uppercase tracking-widest text-gray-500 block mb-1.5">
        Month
      </label>
      <select
        value={selectedMonth ?? ''}
        onChange={e => onMonthChange(e.target.value)}
        disabled={monthsLoading || !months.length}
        className={selectClass}
      >
        {months.map(month => (
          <option key={month} value={month}>{formatMonth(month)}</option>
        ))}
      </select>
    </div>
  </div>
)

export default RidershipFilters
