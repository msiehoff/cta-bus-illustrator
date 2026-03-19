import type { RidershipType } from '../types/api'

interface Props {
  availableMonths: string[]
  selectedMonth: string | null
  ridershipType: RidershipType
  onMonthChange: (month: string) => void
  onTypeChange: (type: RidershipType) => void
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

const selectClass = 'w-full bg-gray-800 border border-gray-700 text-white text-sm rounded-lg px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:opacity-50 cursor-pointer'

const FilterBar = ({ availableMonths, selectedMonth, ridershipType, onMonthChange, onTypeChange }: Props) => (
  <div className="absolute top-3 left-3 z-10 pointer-events-auto">
    <div className="bg-gray-900/90 backdrop-blur border border-gray-700/60 rounded-xl shadow-2xl p-3 flex flex-col gap-3 w-44">
      <div>
        <label className="text-gray-400 text-xs uppercase tracking-widest block mb-1.5">Ridership</label>
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
        <label className="text-gray-400 text-xs uppercase tracking-widest block mb-1.5">Month</label>
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
    </div>
  </div>
)

export default FilterBar
