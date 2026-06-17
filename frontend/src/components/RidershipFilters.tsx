import type { RidershipType } from '../types/api'
import MonthSelector from './MonthSelector'

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
  <div className="flex flex-wrap items-end gap-4 mb-5">
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
    <MonthSelector
      months={months}
      selectedMonth={selectedMonth}
      onMonthChange={onMonthChange}
      disabled={monthsLoading || !months.length}
    />
  </div>
)

export default RidershipFilters
