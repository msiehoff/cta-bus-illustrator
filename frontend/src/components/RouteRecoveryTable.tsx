import {
  formatMonth,
  formatPct,
  formatRecoveryPct,
  formatRides,
  recoveryColorClass,
} from '../lib/ridershipUtils'
import type { RecoveryRow } from '../lib/ridershipUtils'

interface Props {
  currentMonth: string
  benchmarkMonth: string
  yearAgoMonth: string
  fiveYearsAgoMonth: string
  rows: RecoveryRow[]
}

const formatDelta = (pct: number | null) => {
  if (pct == null) return <span className="text-gray-600">—</span>
  return (
    <span className={pct >= 0 ? 'text-green-400' : 'text-red-400'}>
      {formatPct(pct)}
    </span>
  )
}

const RouteRecoveryTable = ({
  currentMonth,
  benchmarkMonth,
  yearAgoMonth,
  fiveYearsAgoMonth,
  rows,
}: Props) => (
  <div className="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden mb-5">
    <div className="px-4 sm:px-5 py-3 border-b border-gray-800">
      <h2 className="text-sm font-medium text-white">Ridership comparison</h2>
      <p className="text-xs text-gray-500 mt-0.5">
        Same calendar month comparisons · {formatMonth(currentMonth)} · pre-pandemic = same month in 2019
      </p>
    </div>
    <div className="overflow-x-auto">
      <table className="w-full text-sm min-w-[28rem]">
        <thead>
          <tr className="bg-gray-800/50 text-xs text-gray-500">
            <th className="text-left px-5 py-2 font-normal">Day type</th>
            <th className="text-right px-3 py-2 font-normal">
              Now
              <span className="block text-[10px] text-gray-600 font-normal">{formatMonth(currentMonth)}</span>
            </th>
            <th className="text-right px-3 py-2 font-normal" title="Same month in 2019, before the pandemic">
              vs pre-pandemic
              <span className="block text-[10px] text-gray-600 font-normal">{formatMonth(benchmarkMonth)}</span>
            </th>
            <th className="text-right px-3 py-2 font-normal hidden sm:table-cell" title="Change vs same month one year ago">
              1y change
              <span className="block text-[10px] text-gray-600 font-normal">{formatMonth(yearAgoMonth)}</span>
            </th>
            <th className="text-right px-5 py-2 font-normal hidden sm:table-cell" title="Change vs same month five years ago">
              5y change
              <span className="block text-[10px] text-gray-600 font-normal">{formatMonth(fiveYearsAgoMonth)}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          {rows.map(row => (
            <tr key={row.type} className="border-t border-gray-800">
              <td className="px-5 py-2.5 text-white">{row.label}</td>
              <td className="px-3 py-2.5 text-right text-gray-300">{formatRides(row.current)}</td>
              <td className="px-3 py-2.5 text-right">
                {row.recoveryPct != null ? (
                  <span className={recoveryColorClass(row.recoveryPct)}>
                    {formatRecoveryPct(row.recoveryPct)}
                  </span>
                ) : (
                  <span className="text-gray-600">—</span>
                )}
              </td>
              <td className="px-3 py-2.5 text-right hidden sm:table-cell">
                {formatDelta(row.yearAgoPct)}
              </td>
              <td className="px-5 py-2.5 text-right hidden sm:table-cell">
                {formatDelta(row.fiveYearPct)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
)

export default RouteRecoveryTable
